package sign

import (
    `context`

    `github.com/chromedp/chromedp`
    log `github.com/sirupsen/logrus`
    `github.com/storezhang/gos/chromedps`

    `songjiang/utils`
)

// G4KQun 4KQun对象
type G4KQun struct {
    HomeUrl        string `default:"http://www.4kqun.com" yaml:"homeUrl" toml:"homeUrl"`
    SignSelector   string `default:"'#JD_sign'" yaml:"signSelector" toml:"signSelector"`
    SignedSelector string `default:"//h1[contains(text(), '您今天已经签到过了或者签到时间还未开始')]" yaml:"signedSelector" toml:"signedSelector"`
    SignUrl        string `default:"http://www.4kqun.com/plugin.php?id=dsu_paulsign:sign" yaml:"signUrl" toml:"signUrl"`
    ScoreUrl       string `default:"http://www.4kqun.com/home.php?mod=spacecp&ac=credit&showcredit=1" yaml:"scoreUrl" toml:"scoreUrl"`
    JBSelector     string `default:"//em[contains(text(), '金币')]/.." yaml:"jbSelector" toml:"jbSelector"`
}

// AutoSign G4KQun的自动签到任务
func (g4kQun *G4KQun) AutoSign(ctx context.Context, cookies string) (result AutoSignResult, err error) {
    // 进入主页
    if e := chromedp.Run(
        ctx,
        chromedps.DefaultVisit(g4kQun.HomeUrl, cookies),
    ); nil != e {
        err = e
        log.WithFields(log.Fields{
            "error": e,
        }).Error("无法载入签到界面")

        return
    } else {
        log.Info("成功进入签到界面")
    }

    // 签到前的K币
    result.Before = getJB(ctx, g4kQun)
    // 确认是否已经签到
    if e := chromedp.Run(
        ctx,
        chromedps.DefaultSleep(),
        chromedps.TasksWithTimeOut(&ctx, "10s", chromedp.Tasks{
            chromedp.Navigate(g4kQun.SignUrl),
            chromedp.WaitVisible(g4kQun.SignedSelector),
        }),
    ); nil != e {
        log.Info("还没有签到，继续执行自动签到任务")
    } else {
        // 签到后的K币
        result.Success = true
        result.After = result.Before
        result.Msg = "已签到，明天再来签到吧"

        log.WithFields(log.Fields{
            "cookies": cookies,
        }).Info("已签到，明天再来签到吧")

        return
    }

    // 点击签到按扭
    if e := chromedp.Run(
        ctx,
        chromedps.DefaultSleep(),
        chromedp.Navigate(g4kQun.SignUrl),
        chromedps.DefaultSleep(),
        chromedp.Click(g4kQun.SignSelector, chromedp.NodeVisible),
    ); nil != e {
        err = e
        log.WithFields(log.Fields{
            "error": e,
        }).Error("无法点击签到按扭")

        return
    } else {
        log.Info("成功点击签到按扭")
        // 签到后的K币
        result.Success = true
        result.After = getJB(ctx, g4kQun)
        result.Msg = "签到成功"
    }

    return
}

func getJB(ctx context.Context, hao4k *G4KQun) (jb string) {
    if err := chromedp.Run(
        ctx,
        utils.Sleep(),
        chromedp.Navigate(hao4k.ScoreUrl),
        chromedp.Text(hao4k.JBSelector, &jb),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法获得当前K币数据")
    } else {
        log.WithFields(log.Fields{
            "currentKB": jb,
        }).Info("成功获得当前K币数据")
    }

    return
}
