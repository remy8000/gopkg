package alert

import (
	"fmt"
	"net/http"
	"net/url"
)

func SendFreeSMS(user, pwd, message string) error {
	apiURL := "https://smsapi.free-mobile.fr/sendmsg"
	params := url.Values{}
	params.Set("user", user)
	params.Set("pass", pwd)
	params.Set("msg", message)

	resp, err := http.Get(fmt.Sprintf("%s?%s", apiURL, params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var message string
		switch resp.StatusCode {
		case 400:
			message = "One of the required parameters is missing"
		case 402:
			message = "Too many SMS have been sent in too short a time"
		case 403:
			message = "The service is not activated on the Subscriber Area, or login/key incorrect"
		case 500:
			message = "Server error. Please try again later"
		default:
			message = "Unknown error"
		}
		return fmt.Errorf("failed to send SMS (code %d): %s", resp.StatusCode, message)
	}
	return nil
}

/*
   200 : Le SMS a été envoyé sur votre mobile.
   400 : Un des paramètres obligatoires est manquant.
   402 : Trop de SMS ont été envoyés en trop peu de temps.
   403 : Le service n'est pas activé sur l'Espace Abonné, ou login / clé incorrect.
   500 : Erreur côté serveur. Veuillez réessayer ultérieurement.
*/
