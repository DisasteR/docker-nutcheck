package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	nut "github.com/robbiet480/go.nut"
	log "github.com/sirupsen/logrus"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Ups :
type Ups struct {
	id      int
	name    string
	model   string
	battery int64
}

var (
	// versionflag : Flag for display version
	versionflag bool
	// debugLevel : Flag for enable debug level
	debugLevel bool
	// version : Default version
	version = "N/A"
	// poolingInterval : pooling interval in minutes
	poolingInterval = 30

	// nutAddress : server address
	nutAddress = "127.0.0.1"
	// nutPort : server port
	nutPort = 3493
	// nutUser : server user
	nutUser string
	// nutPassword : server password
	nutPassword string
	// nutUpsName : ups to check
	nutUpsName string
	// nutUpsList : list ups
	nutUpsList bool

	// telegramToken : To create a bot, please contact @BotFather on telegram
	telegramToken = "None"
	// telegramID : To find an id, please contact @myidbot on telegram
	telegramID = 0

	// Icons
	icon = map[int]string{
		0: "\u2705",       // ":white_check_mark:",
		1: "\u26a0\ufe0f", // ":warning:",
		2: "\u274c",       // ":x:",
		3: "\u2754",       // ":grey_question:",
	}
)

func getUpsVar(ups *nut.UPS, variable string) interface{} {
	for _, upsvar := range ups.Variables {
		if upsvar.Name == variable {
			return upsvar.Value
		}
	}
	return nil
}

func init() {
	// Global
	flag.BoolVar(&versionflag, "v", false, "Print build id")
	flag.BoolVar(&debugLevel, "d", false, "debug mode")

	flag.IntVar(&poolingInterval, "poolinginterval", getIntEnv("POOLING_INTERVAL", poolingInterval),
		"Pooling Interval (or use env variable : POOLING_INTERVAL)")

	// nut Args
	flag.StringVar(&nutAddress, "nutaddress", getStringEnv("NUT_ADDRESS", nutAddress),
		"Nut server address (or use env variable : NUT_ADDRESS)")
	flag.IntVar(&nutPort, "nutport", getIntEnv("NUT_PORT", nutPort),
		"Nut server port (or use env variable : NUT_PORT)")
	flag.StringVar(&nutUser, "nutuser", getStringEnv("NUT_USER", nutUser),
		"Nut server user (or use env variable : NUT_USER)")
	flag.StringVar(&nutPassword, "nutpassword", getStringEnv("NUT_PASSWORD", nutPassword),
		"Nut server password (or use env variable : NUT_PASSWORD)")
	flag.StringVar(&nutUpsName, "nutupsname", getStringEnv("NUT_UPS_NAME", nutUpsName),
		"Nut ups to check (or use env variable : NUT_UPS_NAME)")

	flag.BoolVar(&nutUpsList, "l", false, "List Nut Ups")

	// Telegram
	flag.StringVar(&telegramToken, "telegramtoken", getStringEnv("TELEGRAM_TOKEN", telegramToken),
		"To create a bot, please contact @BotFather on telegram (or use env variable : TELEGRAM_TOKEN)")
	flag.IntVar(&telegramID, "telegramid", getIntEnv("TELEGRAM_ID", telegramID),
		"To find an id, please contact @myidbot on telegram (or use env variable : TELEGRAM_ID)")

	flag.Parse()

	log.SetOutput(os.Stdout)

}

func main() {
	log.Info("Starting...")
	if debugLevel {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatal("Cannot instantiate Telegram BOT")
	}
	if debugLevel {
		bot.Debug = true
	}

	laststate := -1
	for {
		client, err := nut.Connect(nutAddress, nutPort)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Fatal("Exec error")
		}
		_, err = client.Authenticate(nutUser, nutPassword)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Fatal("Exec error")
		}
		upss, err := client.GetUPSList()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Fatal("Exec error")
		}
		client.Disconnect()

		if nutUpsList {
			for idx, ups := range upss {
				fmt.Printf("%+v\n", Ups{
					id:      idx,
					name:    ups.Name,
					model:   getUpsVar(&ups, "device.model").(string),
					battery: getUpsVar(&ups, "battery.charge").(int64),
				})
			}
			return
		}

		for idx, ups := range upss {
			if ups.Name == nutUpsName {
				myups := Ups{
					id:      idx,
					model:   getUpsVar(&ups, "device.model").(string),
					battery: getUpsVar(&ups, "battery.charge").(int64),
				}
				status := 0
				if myups.battery < 20 {
					status = 2
				} else if myups.battery < 50 {
					status = 1
				} else {
					status = 0
				}

				if status != laststate {
					log.Debug("State Transition")

					msg := tgbotapi.NewMessage(int64(telegramID), fmt.Sprintf("%s UPS Model \"%s\" Battery level %d %%", getIcon(status), myups.model, myups.battery))
					_, err = bot.Send(msg)
					if err != nil {
						log.WithFields(log.Fields{
							"err": err.Error(),
						}).Debug("Notification error rescheduling in 5 minutes")
						time.Sleep(time.Duration(int(time.Minute) * 5))
						continue
					}
					log.Debug("Notification sent")
					laststate = status
				}
			} else {
				log.WithFields(log.Fields{
					"nutupsname": nutUpsName,
				}).Fatal("Ups not Found")
			}
		}

		// log.WithFields(log.Fields{
		// 	"args": args,
		// }).Debug("Running command")
		// cmd := exec.Command("check_adaptec_raid", args...)
		// var waitStatus syscall.WaitStatus
		//
		// out, err := cmd.CombinedOutput()
		// if err != nil {
		// 	if exitError, ok := err.(*exec.ExitError); ok {
		// 		waitStatus = exitError.Sys().(syscall.WaitStatus)
		//
		// 		log.WithFields(log.Fields{
		// 			"exitcode": waitStatus.ExitStatus(),
		// 			"out":      out,
		// 		}).Debug("Statuscheck error")
		// 	} else {
		// 		log.WithFields(log.Fields{
		// 			"err": err.Error(),
		// 		}).Fatal("Exec error")
		// 	}
		// } else {
		// 	// Success
		// 	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		// 	log.WithFields(log.Fields{
		// 		"exitcode": waitStatus.ExitStatus(),
		// 		"out":      out,
		// 	}).Debug("Statuscheck error")
		// }
		//
		// if waitStatus.ExitStatus() != laststate {
		// 	log.Debug("State Transition")
		//
		// 	msg := tgbotapi.NewMessage(int64(telegramID), getIcon(waitStatus.ExitStatus())+string(out[:]))
		// 	_, err = bot.Send(msg)
		// 	if err != nil {
		// 		log.WithFields(log.Fields{
		// 			"err": err.Error(),
		// 		}).Debug("Notification error rescheduling in 5 minutes")
		// 		time.Sleep(time.Duration(int(time.Minute) * 5))
		// 		continue
		// 	}
		// 	log.Debug("Notification sent")
		// 	laststate = waitStatus.ExitStatus()
		// }

		log.Debugf("Waiting %d minutes", poolingInterval)
		time.Sleep(time.Duration(int(time.Minute) * poolingInterval))
	}
}

func getIcon(exitcode int) string {
	if i, ok := icon[exitcode]; ok {
		return i + " "
	}
	return ""
}

func getStringEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.WithFields(log.Fields{"key": key}).Info("[main] : Use custom value")
		return value
	}
	log.WithFields(log.Fields{"key": key}).Info("[main] : Use default value")
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.WithFields(log.Fields{"key": key, "err": err}).Fatal("[main] : Invalid value")
			return fallback
		}
		log.WithFields(log.Fields{"key": key}).Info("[main] : Use custom value")
		return int(v)
	}
	log.WithFields(log.Fields{"key": key}).Info("[main] : Use default value")
	return fallback
}
