package sign

import (
    `context`
    `strings`

    `github.com/chromedp/chromedp`
    log `github.com/sirupsen/logrus`
    `github.com/storezhang/gos/chromedps`

    `songjiang/utils`
)

// Hao4k Hao4k对象
type Hao4k struct {
    SignSelector   string `default:"'#JD_sign'" yaml:"signSelector" toml:"signSelector"`
    SignedSelector string `default:"//span[contains(@class, 'btn btnvisted')]" yaml:"signedSelector" toml:"signedSelector"`
    SignUrl        string `default:"https://www.hao4k.cn/k_misign-sign.html" yaml:"signUrl" toml:"signUrl"`
    ScoreUrl       string `default:"https://www.hao4k.cn/home.php?mod=spacecp&ac=credit&showcredit=1" yaml:"scoreUrl" toml:"scoreUrl"`
    KBSelector     string `default:"//em[contains(text(), 'K币')]/.." yaml:"kbSelector" toml:"kbSelector"`
}

// AutoSign Hao4K的自动签到任务
func (hao4k *Hao4k) AutoSign(ctx context.Context, cookies string) (result AutoSignResult, err error) {
    // 等待签到界面完成
    if e := chromedp.Run(
        ctx,
        chromedps.DefaultVisit(hao4k.SignUrl, cookies),
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
    result.Before = strings.TrimSpace(getKB(ctx, hao4k))
    // 确认是否已经签到
    if e := chromedp.Run(
        ctx,
        chromedps.DefaultSleep(),
        chromedps.TasksWithTimeOut(&ctx, "10s", chromedp.Tasks{
            chromedp.Navigate(hao4k.SignUrl),
            chromedp.WaitVisible(hao4k.SignedSelector),
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
        chromedp.Navigate(hao4k.SignUrl),
        chromedps.DefaultSleep(),
        chromedp.Click(hao4k.SignSelector, chromedp.NodeVisible),
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
        result.After = strings.TrimSpace(getKB(ctx, hao4k))
        result.Msg = "签到成功"
    }

    return
}

func getKB(ctx context.Context, hao4k *Hao4k) (kb string) {
    if err := chromedp.Run(
        ctx,
        utils.Sleep(),
        chromedp.Navigate(hao4k.ScoreUrl),
        chromedp.Text(hao4k.KBSelector, &kb),
    ); nil != err {
        log.WithFields(log.Fields{
            "err": err,
        }).Error("无法获得当前K币数据")
    } else {
        log.WithFields(log.Fields{
            "currentKB": kb,
        }).Info("成功获得当前K币数据")
    }

    return
}
