package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
)

var string_s []string

const ( 
    Bot_Token = "PegadinhaDoMalanadro"//(Outro)
    Bot_V     = " v0.3"
    Bot_Name  = "@GoLangCodingBot"
)

func error_check(log_error error){
    if log_error != nil {
        log.Panic(log_error)
    }
}

func Inicia_Bot()(*tgbotapi.BotAPI, error){
    bot, err := tgbotapi.NewBotAPI(Bot_Token)
    error_check(err)

    bot.Debug = true

    return bot, err
}

func Mandar_Mensagen(ChatID int64,Mensagem string){
    bot,err := Inicia_Bot()
    error_check(err)

    msg := tgbotapi.NewMessage(ChatID,Mensagem)
    msg.ParseMode = "html"

    bot.Send(msg)
}

func Responder_Mensagens(ChatID int64,Mensagem string,MensagemID int){
    bot,err := Inicia_Bot()
    error_check(err)

    msg := tgbotapi.NewMessage(ChatID,Mensagem)
    msg.ReplyToMessageID = MensagemID
    msg.ParseMode = "html"

    bot.Send(msg)
}
func Verifica_Comando(mensagem string)bool{
    if string(mensagem[0]) == "/"{
        log.Printf("oi")
        return true
    }else{
        return false
    }
    return false
}
func Comandos(mensagem string) string {
    if strings.Contains(mensagem,"/help"){
        return "<b>Comandos:</b>\n\n/help\n/admins\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }
    if strings.Contains(mensagem,"/admins"){
        return "<b>Caso ocorra algum problema ,não fale com sua mãe, fale com :</b>\n\n@ReiGel_ado<b> ou </b>@Barionix\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }
    return "Comando não encontrado,use o /help para saber os comandos!"
}


func main() {
    bot,err := Inicia_Bot()
    error_check(err)

    log.Printf("Logado na conta  %s", bot.Self.UserName)
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 2 //Tempo pra atualizar

    updates, err := bot.GetUpdatesChan(u)

    for msg := range updates { //Recebe as mensagem
        if msg.Message == nil {
            continue
        }
        
        log.Printf("[%s] %s", msg.Message.From.UserName, msg.Message.Text)
        
        if msg.Message.LeftChatMember != nil{ //Usuario Saiu
            log.Printf("[%s] foi removido do grupo.",msg.Message.LeftChatMember)
            Mandar_Mensagen(msg.Message.Chat.ID,"<b>Esse deve ter feito merda....( ͡° ͜ʖ ͡°)</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>")
        }
        if msg.Message.NewChatMember != nil{
            log.Printf("[%s] foi adicionado/convidado ao grupo.")
            Mandar_Mensagen(msg.Message.Chat.ID,"<b>Eai GOleiro , seja bem vindo a alcateia! Mas conta ai , como chegou aqui ?</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>")
        }

        if msg.Message.Text == ""{
            continue
        }else{
            if Verifica_Comando(msg.Message.Text) == true{
                Mandar_Mensagen(msg.Message.Chat.ID,Comandos(msg.Message.Text))
            }
        }
    }
}
