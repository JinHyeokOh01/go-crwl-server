//정렬 조건 : 날짜순(내림차순)
package crwl

import(
	"time"
	"github.com/JinHyeokOh01/go-crwl-server/models"
)

type NoticeSlice []models.Notice

func (n NoticeSlice) Len() int        { return len(n) }
func (n NoticeSlice) Swap(i, j int)   { n[i], n[j] = n[j], n[i] }
func (n NoticeSlice) Less(i, j int) bool {
    date1, _ := time.Parse("2006-01-02", n[i].Date)
    date2, _ := time.Parse("2006-01-02", n[j].Date)
    return date1.After(date2)
}