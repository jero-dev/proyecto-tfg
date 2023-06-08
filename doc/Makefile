install:
	apt install texlive-full
	apt install openjdk-11-jre-headless
	wget https://github.com/sylvainhalle/textidote/releases/download/v0.8.3/textidote_0.8.3_all.deb -O textidote.deb
	apt install ./textidote.deb
	rm textidote.deb

build:
	pdflatex proyecto.tex && bibtex proyecto && pdflatex proyecto.tex && pdflatex proyecto.tex