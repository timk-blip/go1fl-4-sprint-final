package spentcalories

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	sep := ","
	parts := strings.SplitN(data, sep, 3)
	if len(parts) == 1 && parts[0] == data {
		return 0, "", 0, errors.New("ошибка, сплит не сработал, нет разделителя")
	}
	//if len(parts) < 3 {
	if len(parts) != 3 {
		return 0, "", 0, errors.New("ошибка, неверный формат")
	}
	step := parts[0]
	activity := parts[1]
	if len(parts) == 3 {
		var stepsCurrent int
		stepsCurrent, err := strconv.Atoi(step)
		if err != nil {
			return 0, "", 0, err
		}
		if stepsCurrent <= 0 {
			return 0, "", 0, errors.New("ошибка, 0 шагов")
		}
		durationWalk, err := time.ParseDuration(parts[2])
		if err != nil {
			return 0, "", 0, err
		}
		if durationWalk <= 0 {
			return 0, "", 0, errors.New("ошибка, нулевое время")
		}
		return stepsCurrent, activity, durationWalk, nil
	} 
	return 0, "", 0, nil
}

func distance(steps int, height float64) float64 {
	return (height * stepLengthCoefficient * float64(steps))/float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := time.Duration.Hours(duration)
	speedAverage := dist / hours
	return speedAverage
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsCurrent, activity, durationWalk, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	strings.ToLower(activity)
	switch activity {
	case "Ходьба":
		speed := meanSpeed(stepsCurrent, height, durationWalk)
		str_speed := fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
		dist := distance(stepsCurrent, height)
		str_dist := fmt.Sprintf("Дистанция: %.2f км.\n", dist)
		calories, err := WalkingSpentCalories(stepsCurrent, weight, height, durationWalk)
		if err != nil {
			log.Println(err)
			return "", err
		}
		str_calories := fmt.Sprintf("Сожгли калорий: %.2f\n", calories)
		type_training := fmt.Sprintf("Тип тренировки: %s\n", activity)
		duration := fmt.Sprintf("Длительность: %.2f ч.\n", time.Duration.Hours(durationWalk))
		return type_training + duration + str_dist + str_speed + str_calories, nil
	case "Бег":
		speed := meanSpeed(stepsCurrent, height, durationWalk)
		str_speed := fmt.Sprintf("Скорость: %.2f км/ч\n", speed)
		dist := distance(stepsCurrent, height)
		str_dist := fmt.Sprintf("Дистанция: %.2f км.\n", dist)
		calories, err := RunningSpentCalories(stepsCurrent, weight, height, durationWalk)
		if err != nil {
			log.Println(err)
			return "", err
		}
		str_calories := fmt.Sprintf("Сожгли калорий: %.2f\n", calories)
		type_training := fmt.Sprintf("Тип тренировки: %s\n", activity)
		duration := fmt.Sprintf("Длительность: %.2f ч.\n", time.Duration.Hours(durationWalk))
		return type_training + duration + str_dist + str_speed + str_calories, nil		
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0.0 || height <= 0.0 || duration <= 0 {
		return 0.0, errors.New("некорректные значения")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	minutes := time.Duration.Minutes(duration)
	runSpentCal:= (weight*averageSpeed*minutes)/ minInH
	return runSpentCal, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0.0 || height <= 0.0 || duration <= 0 {
		return 0.0, errors.New("некорректные значения")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	minutes := time.Duration.Minutes(duration)
	runSpentCal:= (weight*averageSpeed*minutes)/ minInH
	correctionСoefficient := runSpentCal * walkingCaloriesCoefficient
	return correctionСoefficient, nil
}
