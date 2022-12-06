package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"time"
)

var document interface{}

var collection *mongo.Collection
var ctx = context.TODO()

type Word struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Eng     string             `bson:"eng,omitempty"`
	Rus     string             `bson:"rus,omitempty"`
	Points  int                `bson:"points,omitempty"`
	Learned bool               `bson:"learned,omitempty"`
}

func getAll() []Word {
	var mapWord []Word
	newWord := collectNewWord()
	midWord := collectMidWord()
	oldWord := collectOldWord()

	mapWord = newWord
	for _, v := range midWord {
		mapWord = append(mapWord, v)
	}
	for _, v := range oldWord {
		mapWord = append(mapWord, v)
	}
	return mapWord

}
func collectOldWord() []Word {
	var end int
	oldWordCursor, err := collection.Find(ctx, bson.D{{"points", bson.M{"$gt": 70}}, {"learned", bson.M{"$eq": false}}})
	if err != nil {
		log.Fatal(err)
	}
	var oldWord []Word
	if err = oldWordCursor.All(ctx, &oldWord); err != nil {
		log.Fatal(err)
	}
	if oldWord == nil {
		return oldWord
	}

	if len(oldWord) > 3 {
		end = 3
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < len(oldWord)/2; i++ {

			ind1 := rand.Intn(len(oldWord))
			ind2 := rand.Intn(len(oldWord))

			oldWord[ind1], oldWord[ind2] = oldWord[ind2], oldWord[ind1]

		}

	} else {
		end = len(oldWord)
	}
	if end <= 0 {
		return oldWord
	}
	oldWord = oldWord[:end]
	return oldWord
}
func collectMidWord() []Word {
	var end int
	midWordCursor, err := collection.Find(ctx, bson.D{{"points", bson.M{"$gt": 30}}, {"points", bson.M{"$lte": 70}}})
	if err != nil {
		log.Fatal(err)
	}
	var midWord []Word
	if err = midWordCursor.All(ctx, &midWord); err != nil {
		log.Fatal(err)
	}

	if midWord == nil {
		return midWord
	}

	if len(midWord) > 5 {
		end = 5
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < len(midWord)/2; i++ {

			ind1 := rand.Intn(len(midWord))
			ind2 := rand.Intn(len(midWord))

			midWord[ind1], midWord[ind2] = midWord[ind2], midWord[ind1]

		}

	} else {
		end = len(midWord)
	}
	if end <= 0 {
		return midWord
	}
	midWord = midWord[:end]
	return midWord
}
func collectNewWord() []Word {
	var end int
	newWordCursor, err := collection.Find(ctx, bson.D{{"points", bson.M{"$lte": 30}}, {"learned", bson.M{"$eq": false}}})
	if err != nil {
		log.Fatal(err)
	}
	var newWord []Word
	if err = newWordCursor.All(ctx, &newWord); err != nil {
		log.Fatal(err)
	}
	if newWord == nil {
		return newWord
	}

	if len(newWord) > 10 {
		end = 10
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < len(newWord)/2; i++ {

			ind1 := rand.Intn(len(newWord))
			ind2 := rand.Intn(len(newWord))

			newWord[ind1], newWord[ind2] = newWord[ind2], newWord[ind1]

		}

	} else {
		end = len(newWord)
	}
	if end <= 0 {
		return newWord
	}
	newWord = newWord[:end]
	return newWord
}
func checkTranslate(wordForCheck []Word) {

	for _, v := range wordForCheck {
		fmt.Print("\nВведите перевод слова: ")
		fmt.Println(v.Eng)
		var rus string
		fmt.Scan(&rus)

		if v.Rus == rus {
			fmt.Print("Перевод правильный. ")
			v.Points += 5
			if v.Points >= 100 {
				v.Learned = true
				result, err := collection.UpdateOne(ctx, bson.M{"_id": v.Id}, bson.D{{"$set", bson.D{{"points", v.Points}, {"learned", v.Learned}}}})
				if err != nil {
					log.Fatal(err)
				}
				log.Println(result.ModifiedCount)
				continue
			}
			result, err := collection.UpdateOne(ctx, bson.M{"_id": v.Id}, bson.D{{"$set", bson.D{{"points", v.Points}}}})
			if err != nil {
				log.Fatal(err)
			}
			log.Println(result.ModifiedCount)

		} else if v.Points > 0 {
			fmt.Println("Перевод неправильный:", v.Eng, " переводится как ", v.Rus)
			v.Points -= 5
			result, err := collection.UpdateOne(ctx, bson.M{"_id": v.Id}, bson.D{{"$set", bson.D{{"points", v.Points}}}})
			if err != nil {
				log.Fatal(err)
			}
			log.Println(result.ModifiedCount)
		} else {
			fmt.Println("Перевод неправильный:", v.Eng, " переводится как ", v.Rus)
		}

	}
	fmt.Println()
}
func addNewWord() {
	var rus string
	var eng string
	var nword int

	fmt.Print("\nВведите количество слов,которые собираетесь ввести: ")
	fmt.Scan(&nword)

	for i := 0; i < nword; i++ {
		fmt.Println("Введите слово на английском языке")
		fmt.Scan(&eng)

		engWordCursor, err := collection.Find(ctx, bson.D{{"eng", bson.M{"$eq": eng}}})
		if err != nil {
			log.Fatal(err)
		}

		if engWordCursor != nil {
			fmt.Println("Слово уже есть в словаре")
			continue
		}

		fmt.Println("Введите слово на русском языке")
		fmt.Scan(&rus)

		document = bson.D{
			{"eng", eng},
			{"rus", rus},
			{"points", 0},
			{"learned", false},
		}
		collection.InsertOne(ctx, document)
	}
	fmt.Println()
}
func deleteWord() {
	var eng string
	fmt.Println("\nВведите слово на английском языке, которое хотите удалить.")
	fmt.Scan(&eng)

	wordForDelete, err := collection.DeleteOne(ctx, bson.D{{"eng", bson.M{"$eq": eng}}})
	if err != nil {
		log.Fatal(err)
	}
	if wordForDelete.DeletedCount > 0 {
		fmt.Printf("Слово %s было удалено\n", eng)
	} else {
		fmt.Printf("Слово %s не найдено в базе\n", eng)
	}
	fmt.Println()
}

func main() {
	clientOptions := options.Client().ApplyURI("")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client.Ping(ctx, nil)

	collection = client.Database("engwords").Collection("words")

	for {
		var choose int
		var all []Word

		fmt.Println("Выберите интересующую вас функцию из списка")
		fmt.Println("1. Добавить слова для изучения")
		fmt.Println("2. Тренировать слова")
		fmt.Println("3. Удалить слово")
		fmt.Println("4. Выход")
		fmt.Print("введите номер пункта меню: ")
		fmt.Scan(&choose)

		if choose == 1 {
			addNewWord()
		} else if choose == 2 {
			all = getAll()
			checkTranslate(all)
		} else if choose == 3 {
			deleteWord()
		} else if choose == 4 {
			break
		} else {
			fmt.Println("Вы ввели некоректный пункт меню:")
		}
	}

}
