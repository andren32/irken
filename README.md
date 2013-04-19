# Irken

Irken är ett pågående arbete med att skapa en grafisk IRC-client i [go](http://golang.org/). För mer information, besök vår [wiki](https://github.com/axelri/irken/wiki) och [projektplan](https://github.com/axelri/irken/wiki/Projektplan).

## Installation

Irken använder sig av ett ramverk som heter GTK, vilket du måste ha installerat för att kunna bygga programmet. Hur man går tillväga varierar mellan olika plattformar.

### Mac OS X

Börja med att installera Xcode med tillhörande kommando-verktyg. GTK använder sig av X, så du måste också installera [XQuartz](http://xquartz.macosforge.org/landing/). Installera sedan [homebrew](http://mxcl.github.io/homebrew/). Nu installerar du lätt go och GTK med:

    brew install go gtk+ gtksourceview glib gdk-pixbuf cairo pango

Följ [instruktionerna för go](http://golang.org/doc/code.html) och sätt upp en GOPATH. Nu installerar du [bindings](https://github.com/mattn/go-gtk/) för go till gtk genom:

    go get github.com/mattn/go-gtk/gtk
    
Får du problem nu kan det vara för att go-gtk inte hittar dina X-filer. Du måste i så fall lägga till följande rad i din ```.bash_profile```:

    export PKG_CONFIG_PATH=/opt/X11/lib/pkgconfig

Nu borde allt vara klart, happy hacking!

### Linux (Testat med Ubuntu 12.10 64-bit)

Installera Go, i Ubuntu kan du med enkelhet göra detta genom att skriva följande i terminalen:

     sudo apt-get install golang

Installera gnome-devel genom att skriva följande i terminalen:

     sudo apt-get install gnome-devel

Följ [instruktionerna för go](http://golang.org/doc/code.html) och sätt upp en GOPATH. Nu installerar du [bindings](https://github.com/mattn/go-gtk/) för go till gtk genom:

    go get github.com/mattn/go-gtk/gtk