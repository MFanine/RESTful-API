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

  subject := "Hey  there!"
  content := `

  <h1>This is an email sent from the Series API server.</h1>
  <h3>Do You know  what is Go?</h3>
   <p>Go is an open-source programming language that makes it easy to build simple, reliable, and efficient software at a local level
 
`

  to := []string{"mouadfanine01@gmail.com"}
  // mohamedouallal02@gmail.com
  attachFiles := []string{"../../README.md"}

  err = sender.SendEmail(subject, content , to, nil, nil ,attachFiles)
  require.NoError(t,err)
  
}
