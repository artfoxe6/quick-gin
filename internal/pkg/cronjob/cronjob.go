package cronjob

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var c *cron.Cron

// JobFormat "syncListing": {"1 2 */1 * *", func() {}},
type JobFormat map[string][2]interface{}

func Schedule(jobs JobFormat) {
	c = cron.New()
	for name, params := range jobs {
		id, err := c.AddFunc(params[0].(string), wrapper(name, params[1].(func())))
		if err != nil {
			log.Printf("[%s] add job err: %s\n", name, err.Error())
		}
		entry := c.Entry(id)

		log.Printf("[%s] job next run time %s", name, entry.Schedule.Next(time.Now()).Format("2006-01-02 15:04:05"))
	}
	c.Start()
}

func wrapper(name string, f func()) func() {
	return func() {
		log.Printf("[%s] job start.\n", name)
		f()
		log.Printf("[%s] job end.\n", name)
	}
}
