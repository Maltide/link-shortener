
## Описание проекта
    URL-shortener allows get a short-version of any link. It uses a hash-function and keeps this pair (original-short) in a database 
    
    We need zap logger

    At first, we need to create a connection with our db, than create a handler with POST-request from users which is send link that we need to short. 

    some handlers for short, redirect, like:

    func hash (s string) string {
        result := 0

        for i := range s {
            result = (result * 7 + int[s[i]]) % 1000
        }

        alphabit := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

        resString := make([]rune, 0, 7)

        for result > 0 && len(resString) != 7 {
            symbolIndex := result % len(alphabit)

            resString = append(resString, rune(alphabit[symbolIndex]))

            result = result / len(alphabit)
        }

        return string(resString)
    }
    func shortHandler(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()

        linkIn := r.FormValue("url")

        shortLink := hash(linkIn)

        add pair in DB (LinkIn == shortLink)

        give user shortversion - need to send html on site where he can see it
        w.Write(shortLink)
    }

    func redirectHandler (w http.ResponseWriter, r *http.Request) {
        r.ParseForm()

        if r.FormValue("url") have in DB{
            redirect on original
        } 
    }
    r := chi.NewRouter()

    in main i need to execute logger, db, handler with methods
    in the end of main i need activate go-server with:

    err := ListenAndServe("rostelecom ip" string, h *http.Handler) 
    if err != nil {
        return log.Error("fail to create server:", err)
    }


