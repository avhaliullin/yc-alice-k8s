// Variables
var messages = document.querySelector('.message-list')

function writeLine(elem, primary) {
    let message = document.createElement('li')
    let actorClass = primary ? 'item-primary' : 'item-secondary'
    message.classList.add('message-item', actorClass)
    message.appendChild(elem)
    messages.appendChild(message)
    messages.scrollTop = messages.scrollHeight;
}

function trimStr(str, limit) {
    if (str.length > limit) {
        str = str.substring(0, limit - 3) + "...";
    }
    return str;
}

function anyField(obj) {
    for (var field in obj) {
        if (obj.hasOwnProperty(field)) {
            return field;
        }
    }
    return null;
}

function hasFields(obj) {
    return anyField(obj) != null;
}

function addMessageProp(msg, key, value) {
    if (msg.hasChildNodes()) {
        msg.appendChild(document.createElement("br"));
    }
    var keyElem = document.createElement("b");
    keyElem.innerText = key + ": ";
    msg.appendChild(keyElem);
    msg.appendChild(document.createTextNode(value));
}

function createElement(tag, ...classes) {
    let elem = document.createElement(tag);
    elem.classList.add(...classes);
    return elem;
}

function renderPropertyLine(key, value, depth) {
    let propElem = createElement("div", "prop-line", "prop-line-" + depth);
    let propKeyElem = createElement("span", "prop-key");
    propKeyElem.innerText = key;
    let propValElem = createElement("span", "prop-val");
    propValElem.innerText = value;
    propElem.appendChild(propKeyElem);
    propElem.appendChild(propValElem);
    return propElem;
}

function renderProperty(root, property, depth) {
    root.appendChild(renderPropertyLine(property.key, property.value, depth));
    if (property["items"]) {
        for (child of property.items) {
            renderProperty(root, child, depth + 1);
        }
    }
}

function renderMessage(messageProps) {
    let root = createElement("div");
    for (const prop of messageProps) {
        renderProperty(root, prop, 0);
    }
    return root
}

function createProperty(key, value) {
    return {key: key, value: value, items: []}
}

function renderRequestMessage(req) {
    const props = [createProperty("request", trimStr(req.request.original_utterance, 100))];
    const intents = req.request.nlu.intents;
    for (const intentName in intents) {
        let renderName = intentName;
        if (renderName.startsWith("easter")) {
            renderName += " " + String.fromCodePoint(0x1F31F);
        }
        const intentProp = createProperty("intent", renderName);
        props.push(intentProp)

        const intentSlots = intents[intentName].slots;
        if (!intentSlots) {
            continue
        }
        for (const slotKey in intentSlots) {
            const slotValue = intentSlots[slotKey].value;
            intentProp.items.push(createProperty(slotKey, slotValue));
        }
    }
    return renderMessage(props);
}

function stateName(stateId) {
    switch (stateId) {
        case "DPLY_REQ_NAME":
            return "Deploying, requested for deploy name";
        case "DPLY_REQ_IMAGE":
            return "Deploying, requested for image name";
        case "DPLY_CNFRM":
            return "Deploying, confirmation requested";
        case "DPLY_ST_REQ_NAME":
            return "Checking deploy status, requested for deploy name";
        case "DPLY_ST_REQ_NMSPC":
            return "Checking deploy status, requested for namespace";
        case "SCL_DPLY_REQ_NAME":
            return "Scaling deploy, requested for deploy name";
        case "SCL_DPLY_REQ_SCALE":
            return "Scaling deploy, requested for new replicas count";
        case "SCL_DPLY_REQ_CNFRM":
            return "Scaling deploy, confirmation requested";
        case "DEL_DPLY_REQ_NAME":
            return "Deleting deploy, deploy name requested";
        case "DEL_DPLY_REQ_CNFRM":
            return "Deleting deploy, confirmation requested";
        case "BRKN_PD_REQ_NS":
            return "Checking broken pods, requested for namespace"
        case "CNT_PD_REQ_NS":
            return "Counting pods, requested for namespace"
        case "LST_INGRS_REQ_NS":
            return "Listing ingresses, requested for namespace"
        case "LST_SRVC_REQ_NS":
            return "Listing services, requested for namespace"
    }
    return stateId;
}

function renderState(state) {
    const stateEnumKey = "State";
    let stateProp = createProperty("state", stateName(state[stateEnumKey]));
    for (const stateKey in state) {
        if (stateKey === stateEnumKey || stateKey.endsWith("ID")) {
            continue;
        }
        stateProp.items.push(createProperty(stateKey, state[stateKey]));
    }
    return stateProp;
}

function renderResponseMessage(resp) {
    const stateEnumKey = "State";
    const props = [createProperty("response", trimStr(resp.response.text, 100))];
    if (resp["session_state"]) {
        props.push(renderState(resp.session_state));
    }
    return renderMessage(props);
}

function onMessage(e) {
    let msg = JSON.parse(e.data);
    let req = msg.req;
    if (req.session["new"]) {
        return;
    }
    let resp = msg.resp;

    writeLine(renderRequestMessage(req), true);
    writeLine(renderResponseMessage(resp), false);

    window.scrollTo(0, document.body.scrollHeight);
}

function connect() {
    var ws = new WebSocket('wss://d5djh2iettflqjen84pm.apigw.yandexcloud.net/ws');

    ws.onmessage = onMessage;

    ws.onclose = function (e) {
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function () {
            connect();
        }, 1000);
    };

    ws.onerror = function (err) {
        console.error('Socket encountered error: ', err.message, 'Closing socket');
        ws.close();
    };
}

connect();