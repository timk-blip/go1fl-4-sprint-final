package daysteps

import (
	//"errors"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	sep := ","
	var durationWalk time.Duration
	var stepsCurrent int
	var err error
	parts := strings.SplitN(data, sep, 2)
	if len(parts) == 1 && parts[0] == data {
		return 0, 0, errors.New("ошибка, сплит не сработал, нет разделителя\n")
	}
	// Нет смысла в этой проверке, так как если мы прошли 35 строку, то точно части 2
	// if len(parts) < 2 ||  len(parts) > 2{
	// 	return 0, 0, errors.New("ошибка, неверный формат\n")
	// }
	// Можно проще, "len(parts) < 2 ||  len(parts) > 2" заменить на !=
	// if len(parts) != 2 {
	// 	return 0, 0, errors.New("ошибка, неверный формат\n")
	// }
	if len(parts) == 2 {
		steps := parts[0] 
		time_walk := parts[1]
		stepsCurrent, err = strconv.Atoi(steps)
		if err != nil {
		return 0, 0, err
		}
		if stepsCurrent <= 0 {
			return 0, 0, errors.New("ошибка, 0 шагов\n")
		}
		durationWalk, err = time.ParseDuration(time_walk)
		if err != nil {
			//return 0, 0, errors.New("ошибка, неверный формат\n")
			return 0, 0, fmt.Errorf("conversion error: %w", err)
		}
		if durationWalk <= 0 {
			return 0, 0, errors.New("ошибка, продолжительность 0\n")
		}
	}
	return stepsCurrent, durationWalk, nil
}

func DayActionInfo(data string, weight, height float64) (string) {
	steps, duration, err :=  parsePackage(data)
	if err != nil {
		log.Println(err)
	//	fmt.Printf("%v", err)
		return ""
	}
	if steps <= 0 {
	//	fmt.Printf("%v", err)
		log.Println(err)
		return ""
	}
	distanceKm := (stepLength * float64(steps)) / float64(mInKm)
	strDist := fmt.Sprintf("Дистанция составила %.2f км.\n", distanceKm)
	strSteps := fmt.Sprintf("Количество шагов: %d.\n", steps)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
	//	fmt.Printf("%v", err)
		log.Println(err)
		return ""
	}
	strCalories := fmt.Sprintf("Вы сожгли %.2f ккал.\n", calories)
	return strSteps + strDist + strCalories
	
}
