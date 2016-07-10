package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
    "io/ioutil"
    "strconv"
)

var(
    ReiGel_ado int = 0
    Barionix   int = 0
    id_string string
)

const ( 
    Bot_Token = ""
    //Bot_Token = ""
    Bot_V     = " v0.4"
    Bot_Name  = "@GoLangCodingBot"
    Rules     = "rules.txt" 
    TDV       = "txt_da_vergonha.txt"
)

func error_check(log_error error){
    if log_error != nil {
        //log.Panic(log_error)
    }
}
/////////////////////////////FUNÇÕES RELACIONADAS AO BOT/SERVIDOR///////////////////////////////////////////
func Inicia_Bot()(*tgbotapi.BotAPI, error){
    bot, err := tgbotapi.NewBotAPI(Bot_Token)
    error_check(err)

    bot.Debug = false

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
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////FUNÇÕES RELACIONADAS A COMANDO/////////////////////////////////
func Verifica_Comando(mensagem string)bool{
    if string(mensagem[0]) == "/"{
        return true
    }else{
        return false
    }
    return false
}
func Comandos(mensagem string,id_usuario int,username string) string {
    ///////////////////////REGRAS//////////////////////////////////Q
    id_string = strconv.Itoa(id_usuario)
    if strings.Contains(mensagem,"/func_regras"){
        novas_regras := strings.Replace(mensagem,"/func_regras","",-1) //Tira o comando da mensagem  :D
        if id_usuario == ReiGel_ado{   
            Escreve(Rules,novas_regras)
            log.Printf("[-]Regras atualizadas por ReiGel_ado!")
            return "<b>As regras foram atualizadas com sucesso pelo admin ReiGel_ado!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }else if id_usuario == Barionix{
            Escreve(Rules,novas_regras)
            log.Printf("[-]Regras atualizadas por Barionix!")
            return "<b>As regras foram atualizadas com sucesso pelo admin Barionix!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }else{
            txt_da_vergonha := Leitor(TDV)
            escreve_txt := txt_da_vergonha + "-" + username + " - ID:" + strconv.Itoa(id_usuario)
            Escreve(TDV,escreve_txt)
            log.Printf("[-]O usuario " + username + " de ID:" + strconv.Itoa(id_usuario) + " tentou alterar as regras!")
            return "<b>Temos um engraçadinho!\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }
    }else if strings.Contains(mensagem,"/regras"){
        regras := Leitor(Rules)
        return regras + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    ////////////////////////////////////////////////////////////////////Q
    }else if strings.Contains(mensagem,"/help"){
        return "<b>Comandos:</b>\n\n/help\n/admins\n/regras\n/func_regras - (Somente Admins)\n/txt_da_vergonha\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }else if strings.Contains(mensagem,"/admins"){
        return "<b>Caso ocorra algum problema ,não fale com sua mãe, fale com :</b>\n\n@ReiGel_ado<b> ou </b>@Barionix\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }else if strings.Contains(mensagem,"/txt_da_vergonha"){
        tdv := Leitor(TDV)
        return "<b>###########MURAL DA VERGONHA###########</b>\n\n" + tdv + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }
    return "Comando não encontrado,use o /help para saber os comandos!"
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////FUNÇÕES RELACIONADAS AO SISTEMA//////////////////////////////////////
func Leitor(arquivo string) string{
    conteudo_arquivo , err := ioutil.ReadFile(arquivo)
    error_check(err)
    return string(conteudo_arquivo)
}
func Escreve(arquivo string,conteudo string){
    conteudo_arquivo := []byte(conteudo)
    err := ioutil.WriteFile(arquivo, conteudo_arquivo , 0777)
    error_check(err)
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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
                Mandar_Mensagen(msg.Message.Chat.ID,Comandos(msg.Message.Text,msg.Message.From.ID,msg.Message.From.UserName))
            }
        }
    }
}
