import pkg from 'ydb-sdk';
import {cloudApi, serviceClients, Session} from '@yandex-cloud/nodejs-sdk';

const {Driver, getLogger, getSACredentialsFromJson, IamAuthService, MetadataAuthService} = pkg;

const {serverless: {apigateway_connection_service: connectionService}} = cloudApi;

let cloudApiSession;
let wsClient;
let wsClientInitialized = false;

async function initWsClient() {
    if (wsClientInitialized) {
        return;
    }
    cloudApiSession = new Session();
    wsClient = cloudApiSession.client(serviceClients.WebSocketConnectionServiceClient);
    wsClientInitialized = true;
}

async function sendMessage(connectionId, message) {
    const request = connectionService.SendToConnectionRequest.fromPartial({
        connectionId,
        type: connectionService.SendToConnectionRequest_DataType.TEXT,
        data: Buffer.from(message, 'utf8'),
    });

    return wsClient.send(request);
}

let logger = getLogger();
let driver;

async function initDb() {
    logger.info('Driver initializing...');

    let authService;
    if (process.env.SA_KEY_FILE) {
        const saKeyFile = process.env.SA_KEY_FILE;
        const saCredentials = getSACredentialsFromJson('./' + saKeyFile);
        authService = new IamAuthService(saCredentials);
    } else {
        authService = new MetadataAuthService();
    }

    driver = new Driver({
        endpoint: process.env.ENDPOINT,
        database: process.env.DATABASE,
        authService,
    });
    const timeout = 10000;
    if (!(await driver.ready(timeout))) {
        logger.fatal(`Driver has not become ready in ${timeout}ms!`);
        process.exit(1);
    }
    logger.info('Driver ready');
}


async function registerConnection(connectionID) {
    const query = `
    UPSERT INTO connections (id)
    VALUES ('${connectionID}');
  `;

    await driver.tableClient.withSession(async (session) => {
        await session.executeQuery(query);
    });
}

async function unregisterConnection(connectionID) {
    const query = `
        DELETE
        FROM connections
        WHERE id == '${connectionID}';
    `;
    await driver.tableClient.withSession(async (session) => {
        return await session.executeQuery(query);
    });
}

async function listConnections() {
    const query = `SELECT id
                   FROM connections;`;
    const {resultSets} = await driver.tableClient.withSession(async (session) => {
        return await session.executeQuery(query);
    });
    let connections = [];
    let rs = resultSets[0];
    for (const row of rs.rows) {
        connections.push(row.items[0].textValue);
    }
    return connections;
}

let dbInitialized = false

async function initDbIfNeeded() {
    if (!dbInitialized) {
        await initDb();
        dbInitialized = true;
    }
}

export const handler = async function (event) {
    logger.info(event);
    await initDbIfNeeded()
    let ctx = event["requestContext"];
    if (ctx["eventType"] === "CONNECT") {
        await registerConnection(ctx.connectionId);
        return {"statusCode": "200"};
    } else if (ctx["eventType"] === "DISCONNECT") {
        await unregisterConnection(ctx.connectionId);
        return {"statusCode": "200"};
    } else {
        return {"statusCode": "200"};
    }
}

async function sendToWs(conns, msg) {
    logger.info("sending event: ", msg)
    for (const connId of conns) {
        try {
            logger.info("sending to ", connId)
            await sendMessage(connId, msg);
        } catch (e) {
            logger.warn("error sending message to " + connId, e);
        }
    }
}

export const onLogEvent = async function (event) {
    logger.info(event);
    await initDbIfNeeded();
    await initWsClient();

    let conns = await listConnections();

    for (const msgs of event.messages) {
        for (const msg of msgs.details.messages) {
            if (!msg["json_payload"]) {
                continue;
            }
            let payload = msg.json_payload;
            if (payload["kind"] !== "access-log") {
                continue;
            }
            let writeMsg = JSON.stringify({req: payload.req, resp: payload.resp});
            await sendToWs(conns, writeMsg)
        }
    }
    return {};
}