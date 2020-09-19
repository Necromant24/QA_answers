package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"io/ioutil"
	"log"
)

type QA struct{
	Contexts []string    `json:"context_raw"`
	Questions []string   `json:"question_raw"`
}




func main() {


	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/answer",func(ctx *gin.Context) {
		body,err:=ctx.Request.GetBody()
		if err!=nil{
			panic(err)
		}
		bytes,err:= ioutil.ReadAll(body)
		if err!=nil{
			panic(err)
		}
		dictData := make(map[string]map[string]string)
		jsonData := make(map[string]interface{})
		err = json.Unmarshal(bytes,&jsonData)
		dictData = jsonData["questions"].(map[string]map[string]string)

		answers:= getMapAnswers(dictData,jsonData["context"].(string))

		jsonStr,err := json.Marshal(answers)
		if err!=nil{
			panic(err)
		}
		ctx.JSON(200,gin.H{
			"data": string(jsonStr),
		})

	})





	r.Run()


}

func test(){
	dat, err := ioutil.ReadFile("datafiles/context.txt")
	if err!=nil{
		panic(err)
	}
	var ctxdat = []string{string(dat)}

	fmt.Println(getAnswer(QA{Contexts:ctxdat,Questions:[]string{"Где работала?"}}))

	btext,err := ioutil.ReadFile("datafiles/ofiziant.json")
	if err!=nil{
		panic(err)
	}
	var jsonObj = make(map[string]map[string]string)
	err = json.Unmarshal(btext,&jsonObj)
	fmt.Println(jsonObj["q1"]["question"])


	fmt.Println("--------")

	answers:= getMapAnswers(jsonObj, string(dat))
	for key,val:=range answers{
		fmt.Println(key+" : "+val)
	}
}





func getMapAnswers(obj map[string]map[string]string, context string)map[string]string{
	var answers = make(map[string]string)
	for key,val:=range obj{
		answers[key]= getAnswer(QA{Contexts:[]string{context},Questions:[]string{val["question"]}})
	}
	return answers
}


func getMapQuestions(obj map[string]map[string]string)map[string]string{
	var questions map[string]string
	for key,val:=range obj{
		questions[key]=val["question"]
	}
	return questions
}



func getAnswer(qa QA)string{
	header := req.Header{
		"Accept":        "application/json",
	}
	param := req.BodyJSON(&qa)

	r, err := req.Post("https://7005.lnsigo.mipt.ru/model", header, param)
	if err != nil {
		log.Fatal(err)
	}

	var s,err2 = ioutil.ReadAll(r.Response().Body)
	if err!=nil{
		panic(err2)
	}

	return string(s)
}

func getAnswerList(questions []string,context string)[]string{
	answers:= []string{}
	ctx:= []string{context}

	for _,val := range questions{
		var qstns = []string{val}
		answers = append(answers,getAnswer(QA{Contexts:ctx,Questions:qstns}))
	}

	return  answers
}


