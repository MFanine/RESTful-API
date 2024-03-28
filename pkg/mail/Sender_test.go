package mail

import (
	"series/pkg/utils"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T){
  config, err := utils.LoadConfig("..")
  require.NoError(t, err)

  sender :=  NewGmailSender(config.EmailSenderName , config.EmailSenderAddress, config.EmailSenderPassword)

  subject := "A test email"
  content := `
  <h1>Hello World!</h1>
  <p>This is a test email.</p>
`

  to := []string{"mouadfanine@gmail.com"}
  attachFiles := []string{"../../README.md"}

  err = sender.SendEmail(subject, content , to, nil, nil ,attachFiles)
  require.NoError(t,err)
  
}
