package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// TEURL is base url for data
	TEURL = "https://paikat.te-palvelut.fi/tpt-api/tyopaikat?valitutAmmattialat=25&valitutAmmattialat=35&ilmoitettuPvm=1&vuokrapaikka=---&sort=mainAmmattiRivino%20asc,%20tehtavanimi%20asc,%20tyonantajanNimi%20asc,%20viimeinenHakupaivamaara%20asc&kentat=ilmoitusnumero,tyokokemusammattikoodi,ammattiLevel3,tehtavanimi,tyokokemusammatti,tyonantajanNimi,kunta,ilmoituspaivamaara,hakuPaattyy,tyoaikatekstiYhdistetty,tyonKestoKoodi,tyonKesto,tyoaika,tyonKestoTekstiYhdistetty,hakemusOsoitetaan,maakunta,maa,hakuTyosuhdetyyppikoodi,hakuTyoaikakoodi,hakuTyonKestoKoodi&rows=100&start=0&ss=true&facet.fkentat=hakuTyoaikakoodi,ammattikoodi,aluehaku,hakuTyonKestoKoodi,hakuTyosuhdetyyppikoodi,oppisopimus&facet.fsort=index&facet.flimit=-1"
	// TEURLJOB is base url for job description
	TEURLJOB = "https://paikat.te-palvelut.fi/tpt-api/tyopaikat/%v?kieli=fi"
)

// TEData is main response from api
type TEData struct {
	Response TEResponse
}

type TEDataJob struct {
	Response TEResponseJob
}

// TEResponse is actual response data from main response
type TEResponse struct {
	Docs []TEDoc
}

type TEResponseJob struct {
	Docs []TEJob
}

// TEDoc is doc response from api
type TEDoc struct {
	ID    int    `json:"ilmoitusnumero"`
	Title string `json:"tehtavanimi"`
}

// TEJob is the job description struct
type TEJob struct {
	Description string `json:"kuvausteksti"`
}

func main() {

	response, err := http.Get(TEURL)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var data TEData
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	var jobDescriptions []TEJob
	for _, doc := range data.Response.Docs {

		response, err := http.Get(fmt.Sprintf(TEURLJOB, doc.ID))
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		var data TEDataJob
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err)
		}

		jobDescriptions = append(jobDescriptions, data.Response.Docs...)
	}

	bytes, err := json.Marshal(jobDescriptions)

	ioutil.WriteFile("descriptions.json", bytes, 0644)
}
