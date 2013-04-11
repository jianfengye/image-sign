package main

import(
    "signer"
    "utility/process"
    "os"
    "log"
)

var (
    dataPath string = ""
)

func init() {
    rootPath, err := process.RootPath()
    if err != nil{
        log.Fatalln(err)
    }
    dataPath = rootPath + "/dat/"
}

func main() {
    // 输入文件
    src, err := os.Open(dataPath + "src.png")
    if err != nil {
        log.Fatalln(err)
    }
    defer src.Close()

    dst, err := os.OpenFile(dataPath + "out.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
    if err != nil {
        log.Fatalln(err)
    }
    defer dst.Close()
    
    signWriter := signer.NewSigner(src, dst, dataPath + "luximr.ttf")
    signWriter.SetStartPoint(-5,-10)
    err = signWriter.Sign("yejianfeng")
    if err != nil {
        log.Fatalln(err)
    }
}