package mail

func SendSuccessConfirmation(email string, category string) bool {
	switch category {
	case "register":
		sendSuccessRegistration(email)
	}
	return true
}

func sendSuccessRegistration(email string) bool {
	// TODO: implement it
	return true
}
