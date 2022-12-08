package main


import (
    "log"
    "os"
    "strconv"
    "strings"
    "time"
    "math/rand"
    "io/ioutil"
    "net/http"
    "net/url"
    "path/filepath"
)


var (
    CURRENT string;
    OTAKU_LINKS []byte
    FILES = make(map[string]string);
    TEXTS = make(map[string]int);
    PORT = os.Getenv("PORT");
)


func init() {
    cur, err := filepath.Abs(".");
    if err != nil {
        log.Fatal(err);
    }

    if len(PORT) == 0 {
        PORT = "80";
    }

    CURRENT = cur;
    FILES["/"] = filepath.Join(CURRENT, "public",  "index.html");
    FILES["/main.css"] = filepath.Join(CURRENT, "public", "main.css");
    FILES["/main.js"] = filepath.Join(CURRENT, "public", "main.js");
    FILES["/vue.min.js"] = filepath.Join(CURRENT, "public", "vue.min.js");

    locales, _ := ioutil.ReadDir(filepath.Join(CURRENT, "texts"));
    for _, locale := range locales {
        name := locale.Name();
        t, _ := ioutil.ReadDir(filepath.Join(CURRENT, "texts", name));
        TEXTS[name] = len(t);
    }

    links, _ := ioutil.ReadDir(filepath.Join(CURRENT, "images", "otaku"));
    var linkList []string;
    for _, link := range links {
        linkList = append(linkList, "/images/otaku/" + link.Name());
    }

    OTAKU_LINKS = []byte(strings.Join(linkList, ","));

    rand.Seed(time.Now().UTC().UnixNano());
}


func main() {
    http.HandleFunc("/", handleSpecialPath);
    http.HandleFunc("/texts/", handleRawPath);
    http.HandleFunc("/locales/", handleRawPath);
    http.HandleFunc("/images/", handleRawPath);
    http.HandleFunc("/importLinks/", handleImportLinks);
    http.HandleFunc("/random_text/", handleRedirectTexts);
    log.Fatal(http.ListenAndServe(":" + PORT, nil));
}


func handleSpecialPath(res http.ResponseWriter, req *http.Request) {
    path := req.URL.Path;
    if _, ok := FILES[path]; !ok { return; }
    http.ServeFile(res, req, FILES[path]);
}


func handleRawPath(res http.ResponseWriter, req *http.Request) {
    path := filepath.Join(CURRENT, req.URL.Path);
    http.ServeFile(res, req, path);
}


func handleImportLinks(res http.ResponseWriter, req *http.Request) {
    res.Write(OTAKU_LINKS);
}


func handleRedirectTexts(res http.ResponseWriter, req *http.Request) {
    link := getRandomTextFile(getLocale(req.URL.RawQuery));
    http.Redirect(res, req, link, 302);
}


func getLocale(rawQuery string) string {
    urlQueryMap, err := url.ParseQuery(rawQuery);
    if err != nil {
        return "en";
    }

    if val, ok := urlQueryMap["locale"]; ok {
        return val[0];
    }

    return "en";
}


func getRandomTextFile(locale string) string {
    textsLen := TEXTS[locale];
    ind := strconv.Itoa(rand.Intn(textsLen));
    return "/texts/" + locale + "/" + ind;
}
