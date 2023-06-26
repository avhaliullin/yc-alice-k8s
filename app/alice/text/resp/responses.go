// Copyright (c) 2023 Yandex LLC. All rights reserved.
// Author: Andrey Khaliullin <avhaliullin@yandex-team.ru>

package resp

import (
	"fmt"
	"math/rand"
	"strings"

	aliceapi "github.com/avhaliullin/yandex-alice-k8s-skill/app/alice/api"
	"github.com/avhaliullin/yandex-alice-k8s-skill/app/k8s"
)

func AskNSForBrokenPods() *aliceapi.Response {
	return format("В каком нэймсп+эйсе поискать сломанные п+оды?")()
}

func AskNSForServiceList() *aliceapi.Response {
	return randomize(
		format("В какой нэймсп+эйсе поискать сервисы?"),
		format("В какой нэймсп+эйсе посмотреть сервисы?"),
		format("Сервисы какого нэймсп+эйса вас интересуют?"),
	)()
}

func NSNotFound(ns string) *aliceapi.Response {
	return randomize(
		format("Я не нашла нэймсп+эйс \"%s\"", ns),
		format("Не смогла найти нэймсп+эйс \"%s\"", ns),
	)()
}

func BrokenPodsInNS(ns string, statuses *k8s.PodStatusesResp) *aliceapi.Response {
	brokenPods := statuses.Failed
	resp := numDependent(brokenPods, numDependentConfig{
		exactly0: format("В нэймсп+эйсе \"%s\" нет сломанных п+одов", ns),
		like1:    format("В нэймсп+эйсе \"%s\" %d сломанный под", ns, brokenPods),
		like2:    format("В нэймсп+эйсе \"%s\" %d сломанных п+ода", ns, brokenPods),
		like5:    format("В нэймсп+эйсе \"%s\" %d сломанных п+одов", ns, brokenPods),
	})
	running := statuses.Succeeded + statuses.Running
	if running > 0 {
		resp = concat(resp,
			numDependent(running,
				numDependentConfig{
					like1: format(", %d п+од запущен", running),
					like2: format(", %d п+ода запущено", running),
					like5: format(", %d п+одов запущено", running),
				},
			))
	}
	if statuses.Pending > 0 {
		resp = concat(resp, format(", %d еще в процессе запуска", statuses.Pending))
	}
	return resp()
}

func AskNSForCountingPods() *aliceapi.Response {
	return format("В каком нэймсп+эйсе посчитать п+оды?")()
}

func PodsCountInNS(ns string, podsCount int) *aliceapi.Response {
	return numDependent(podsCount, numDependentConfig{
		exactly0: format("В нэймсп+эйсе \"%s\" нет п+одов", ns),
		like1:    format("В нэймсп+эйсе \"%s\" %d под", ns, podsCount),
		like2:    format("В нэймсп+эйсе \"%s\" %d п+ода", ns, podsCount),
		like5:    format("В нэймсп+эйсе \"%s\" %d п+одов", ns, podsCount),
	})()
}

func RejectOnWizard() *aliceapi.Response {
	return randomize(
		format("Хорошо, отменяю, давайте попробуем что-нибудь еще"),
		format("Отменяю"),
		format("Окей, отменяю. Чем я могу помочь?"),
	)()
}

func WhichDeployToDelete() *aliceapi.Response {
	return format("Какой депл+ой удалить?")()
}

func DeployNotFound(deploymentName string) *aliceapi.Response {
	return randomize(
		format("Я не нашла депл+оймэнт \"%s\"", deploymentName),
		format("Я не нашла депл+ой \"%s\"", deploymentName),
		format("Не смогла найти депл+ой \"%s\"", deploymentName),
	)()
}

func ConfirmDeletingDeploy(name string) *aliceapi.Response {
	return randomize(
		format("Удаляю депл+ой \"%s\". Все верно?", name),
		format("Удаляю депл+оймент \"%s\". Все верно?", name),
	)()
}

func DeployDeletionFailed(name string) *aliceapi.Response {
	return randomize(
		format("Не получилось удалить депл+ой"),
		format("Ну удалось удалить депл+ой, посмотрите детали в логах"),
	)()
}

func DeployDeleted(name string) *aliceapi.Response {
	return randomize(
		format("Готово, запустила удаление депл+оя"),
		format("Отправила депл+оймент удаляться"),
	)()
}

func DeployScaleMinAssert() *aliceapi.Response {
	return format("Я не могу задепл+оить меньше одной реплики")()
}

func DeployScaleMaxAssert(max int) *aliceapi.Response {
	return numDependent(max, numDependentConfig{
		like1: format("Я могу задепл+оить не больше %d реплики", max),
		like2: format("Я могу задепл+оить не больше %d реплик", max),
		like5: format("Я могу задепл+оить не больше %d реплик", max),
	})()
}

func WhichImageToDeploy() *aliceapi.Response {
	return format("Какой образ мне задепл+оить?")()
}

func ImageNotFound(image string) *aliceapi.Response {
	return randomize(
		format("Я не знаю образ \"%s\"", image),
		format("Не нашла образ \"%s\"", image),
		format("Извините, я не умею депл+оить образ \"%s\"", image),
	)()
}

func HowToNameDeploy() *aliceapi.Response {
	return randomize(
		format("Как мы назовем депл+ой?"),
		format("Как назвать депл+ой?"),
		format("Какое имя дать депл+ою?"),
	)()
}

func ConfirmDeploy(name string, image string, scale int) *aliceapi.Response {
	scaleNumStr := number(scale, CaseAccusative, GenderF)
	replicas := numDependent(scale, numDependentConfig{
		like1: format("%s реплику", scaleNumStr),
		like2: format("%s реплики", scaleNumStr),
		like5: format("%s реплик", scaleNumStr),
	})
	return randomize(
		concat(
			format("Запускаю депл+ой \"%s\" из образа \"%s\" на ", name, image),
			replicas,
			format(". Все верно?"),
		),
		concat(
			format("Вы хотите начать депл+ой \"%s\" из образа \"%s\" на ", name, image),
			replicas,
			format(". Правильно?"),
		),
	)()
}

func DeployFailed() *aliceapi.Response {
	return randomize(
		format("Не получилось запустить депл+ой"),
		format("Я не смогла запустить депл+ой, проверьте мои логи"),
	)()
}

func DeployStarted() *aliceapi.Response {
	return randomize(
		format("Готова, запустила депл+ой"),
		format("Депл+ой поехал"),
		format("Отлично, покатилось!"),
	)()
}

func AskNSForDeployStatus() *aliceapi.Response {
	return format("В каком нэймсп+эйсе проверить депл+ой?")()
}

func DeployNameForStatus(availableNames []string) *aliceapi.Response {
	depListStr := strings.Join(availableNames, "\n")
	return format("Как называется депл+ой? Их тут %d: %s", len(availableNames), depListStr)()
}

func DeployNotFoundInNS(ns string, deployName string) *aliceapi.Response {
	return format("Я не нашла депл+ой \"%s\" в нэймсп+эйсе \"%s\"", deployName, ns)()
}

func DeployScalingConfirm(name string, scale int) *aliceapi.Response {
	return numDependent(scale, numDependentConfig{
		like1: format("Масштабирую депл+ой \"%s\" до %d реплики. Все верно?", name, scale),
		like2: format("Масштабирую депл+ой \"%s\" до %d реплик. Все верно?", name, scale),
		like5: format("Масштабирую депл+ой \"%s\" до %d реплик. Все верно?", name, scale),
	})()
}

func DeployScalingFail(name string) *aliceapi.Response {
	return randomize(
		format("Не получилось отмасштабировать депл+ой"),
		format("Масштабирование не получилось. Попробуйте посмотреть в моих л+огах, что пошло не так"),
	)()
}

func DeployScalingSuccess(ns, name string) *aliceapi.Response {
	return randomize(
		format("Запустила масштабирование. "+
			"Чтобы узнать, как дела - попросите меня рассказать статус депл+оя \"%s\" в нэймсп+эйсе \"%s\"", name, ns),
		format("Масштабирование запущено. Если захотите удалить депл+ой - просто попросите."),
	)()
}

func DeployReplicaStatuses(deploy string, available int, unavailable int) *aliceapi.Response {
	if unavailable > 0 {
		return concat(
			numDependent(available, numDependentConfig{
				exactly0: format("В депл+ое \"%s\" нет доступных реплик, ", deploy),
				like1:    format("В депл+ое \"%s\" %d доступная реплика, ", deploy, available),
				like2:    format("В депл+ое \"%s\" %d доступные реплики, ", deploy, available),
				like5:    format("В депл+ое \"%s\" %d доступных реплик, ", deploy, available),
			}),
			numDependent(unavailable, numDependentConfig{
				like1: format("еще %d реплика в статусе анав+эйлабл", unavailable),
				like2: format("еще %d реплики в статусе анав+эйлабл", unavailable),
				like5: format("еще %d реплик в статусе анав+эйлабл", unavailable),
			}),
		)()
	}
	return numDependent(available, numDependentConfig{
		exactly1: format("Единственная реплика в депл+ое \"%s\" уже доступна", deploy),
		exactly2: format("Обе реплики в депл+ое \"%s\" доступны", deploy),
		like1:    format("%d реплика уже доступна в депл+ое \"%s\"", available, deploy),
		like2:    format("Все %d реплики уже доступны в депл+ое \"%s\"", available, deploy),
		like5:    format("Все %d реплик уже доступны в депл+ое \"%s\"", available, deploy),
	})()
}

func AskNSForIngresses() *aliceapi.Response {
	return randomize(
		format("В каком нэймсп+эйсе перечислить ингр+эссы?"),
		format("В каком нэймсп+эйсе поискать ингр+эссы?"),
	)()
}

func ListIngresses(ns string, ingressList []string) *aliceapi.Response {
	if len(ingressList) == 0 {
		return format("В нэймсп+эйсе \"%s\" нет ингр+эссов", ns)()
	}
	ingressListStr := joinItemsList(ingressList)
	return numDependent(len(ingressListStr), numDependentConfig{
		exactly1: format("В нэймсп+эйсе \"%s\" есть только ингр+эсс %s", ns, ingressListStr),
		like1:    format("В нэймсп+эйсе \"%s\" %d ингр+эсс: %s", ns, len(ingressList), ingressListStr),
		like2:    format("В нэймсп+эйсе \"%s\" %d ингр+эсса: %s", ns, len(ingressList), ingressListStr),
		like5:    format("В нэймсп+эйсе \"%s\" %d ингр+эссов: %s", ns, len(ingressList), ingressListStr),
	})()
}

func ListServices(ns string, serviceList []string) *aliceapi.Response {
	if len(serviceList) == 0 {
		return format("В нэймсп+эйсе \"%s\" нет сервисов", ns)()
	}
	serviceListStr := joinItemsList(serviceList)
	return numDependent(len(serviceListStr), numDependentConfig{
		exactly1: format("В нэймсп+эйсе \"%s\" есть только сервис %s", ns, serviceListStr),
		like1:    format("В нэймсп+эйсе \"%s\" %d сервис: %s", ns, len(serviceList), serviceListStr),
		like2:    format("В нэймсп+эйсе \"%s\" %d сервиса: %s", ns, len(serviceList), serviceListStr),
		like5:    format("В нэймсп+эйсе \"%s\" %d сервисов: %s", ns, len(serviceList), serviceListStr),
	})()
}

func ListNSs(nss []string) *aliceapi.Response {
	nssString := joinItemsList(nss)
	return numDependent(len(nss), numDependentConfig{
		like1: format("Я нашла %d нэймсп+эйс: %s", len(nss), nssString),
		like2: format("Я нашла %d нэймсп+эйса: %s", len(nss), nssString),
		like5: format("Я нашла %d нэймсп+эйсов: %s", len(nss), nssString),
	})()
}

func WhichDeployToScale() *aliceapi.Response {
	return randomize(
		format("Какой депл+ой отск+ейлить?"),
		format("Какой депл+ой вы хотите отмасштабировать?"),
	)()
}

func ExpectedNumber() *aliceapi.Response {
	return format("Я вас не поняла, давайте попробуем заново")()
}

func EasterDBLaunch() *aliceapi.Response {
	return format("Если у вас возникает такой вопрос — то нет")()
}

func EasterHowTo() *aliceapi.Response {
	return format("В оупэнш+ифт это уже работает из коробки")()
}

func EasterWhatIsK8s() *aliceapi.Response {
	return format("Куберн+этис — это пять бинар+ей")()
}

func EasterHowYouMade() *aliceapi.Response {
	return format("Немного бесс+ерверных вычислений и очень, очень, очень много +ифов")()
}

func NoDeploymentsInNS(ns string) *aliceapi.Response {
	return format("В нэймсп+эйсе \"%s\" нет деплойм+энтов")()
}

func discoveryBase() respF {
	return format("Я умею заглядывать в нэймсп+эйсы: " +
		"считать и искать сломанные п+оды, " +
		"могу перечислить сервисы и ингр+эссы в нэймсп+эйсе")
}

func ScenarioDiscovery() *aliceapi.Response {
	return discoveryBase()()
}

func LetsPlayK8S() *aliceapi.Response {
	return concat(
		randomize(
			format("Давайте! "),
			format("Обязательно! "),
		),
		discoveryBase(),
	)()
}

func WelcomePhrase() *aliceapi.Response {
	return randomize(
		format("Давайте разберемся с вашим кубером!"),
		format("Руки чешутся чего-нибудь подепл+оить!"),
	)()
}

func UnrecognizedRequest() *aliceapi.Response {
	return randomize(
		format("Я вас не поняла"),
		format("Извините, не понимаю"),
		format("Я ничего не поняла, но могу рассказать вам, что я умею"),
		format("А какую задачу мы решаем?"),
		format("Непонятно. Так-то я искусственный интеллект, а адм+иню так, для души"),
	)()
}

func joinItemsList(items []string) string {
	trimmed := false
	if len(items) > 5 {
		rand.Shuffle(len(items), func(i, j int) {
			x := items[i]
			items[i] = items[j]
			items[j] = x
		})
		items = items[:5]
		trimmed = true
	}
	result := ""
	for _, item := range items {
		if len(result) > 0 {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\"", item)
	}
	if trimmed {
		result += " и другие"
	}
	return result
}
