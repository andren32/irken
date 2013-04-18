# Projektplan
## Programbeskrivning
Programmet ska kunna koppla upp sig på populära IRC-servrar och kunna stödja grundläggande operationer i IRC-protokollet. Med hjälp av Go:s multitrådning hoppas vi också kunna implementera ett smidigt sätt att samtidigt vara uppkopplad mot flera kanaler och/eller servrar.

## Användarbeskrivning
Vår målgrupp är ganska vana datoranvändare som redan har viss erfarenhet av att använda andra IRC-klienter. Vi kommer dock att sikta på att implementera ett hjälpkommando som kan täcka åtminstone den grundläggande funktionaliteten.

## Användarscenarion

### I
Användaren går in i vårt program och ser två textrutor och en send-knapp. I den översta rutan står det ”Welcome to irken. For help type /help”. Vår användare skriver ”/help” i den nedre textrutan och trycker sedan på ”send”. Han får då upp en lista på kommandon han kan använda. Han väljer att gå in på QuakeNet och chatta på kanalen #kthd. Han skriver därför ”/connect quakenet” och sedan ”/join #kthd”. För att folk ska känna igen honom skriver han ”/nick crazybanana92” och blir då ett med sitt alterego. Han skriver sitt första inlägg ”vim drools, emacs RULES!”, trycker på ”send” och blir sedermera utkastad ur kanalen.

### II
Användaren går tillbaka in i #kthd och tar sig samman. Han vill nu chatta i en annan kanal, men han vill samtidigt ta emot meddelanden från #kthd. Han går upp i en meny och väljer ”ny tabb”. Han hamnar då i ett nytt fönster, men har kvar sin uppkoppling mot QuakeNet och sitt fina nick. Han väljer att gå in på en annan kanal. Han skriver då ”/join #fabhack” och pratar en stund på den nya kanalen. Efter tag vill han tillbaka till sin förra kanal. Han väljer då ”#kthd”-tabben från en meny, och kommer tillbaka till den förra kanalen. Han kan då se vad alla andra har skrivit medan han var borta.

## Testplan
Självklart kommer programmet att testas med både manuella och automatiska tester av oss som programmerar. Vi kommer försöka att skriva så mycket testkod vi kan, men vet inte hur lätt detta blir eftersom stor del av programmet består av nätverks- och grafikkod. Testanvändaren kommer att genomföra våra användarscenarion och ge feedback på hur åtkomliga dessa var att göra.

## Programdesign
(Eftersom vi skriver i Go tolkar vi "klasser" som kodmoduler). Vi har i grova drag tänkt implementera följande modulära delar:

* __Conn__ symboliserar en uppkoppling till en viss IRC-server. De viktigaste metoderna att kommer kunna koppla upp sig mot en server och läsa och skriva från denna. __Conn__ ska även ha metoder för att ta reda på vilka kanaler användaren är inne i, användarens användarnamn på servern, samt att kunna avsluta uppkopplingen mot servern.

* __Session__ symboliserar en användning av programmet. __Session__ kommer att innehålla metoder för få ut standard-nick och andra detaljer om användaren. __Session__ kommer också att kunna lista buffertar för alla öppna kanaler.

* __Buffer__ symboliserar en buffert som tar emot utskrifter. __Buffer__ kommer att höra till en viss kanal eller server-initiering.

* __Parser__ tolkar informationen till och från servern och skickar den till rätt __Buffer__ (vid läsning) eller __Conn__ (vid skrivning).

* __Window__ och __Widget__ är fönster och delfönster. En __Widget__ kommer att höra till en viss __Buffer__.

## Tekniska frågor

* Hur fungerar IRC-protokollet och vilka delar måste stödjas för att få en någorlunda användbar klient? Kanske räcker det med att avgränsa sig.

* Hur kan man på ett smidigt sätt dela upp informations-flödet? Kanaler ska visas i olika fönster, men flödet från servern är konstant och ofiltrerat. Kanaler måste kunna hållas öppna utan att de visas i ett fönster. Någon sorts buffertar som man associerar med ett fönster (som i t.ex vim) kanske är en bra lösning.

## Arbetsplanen

* Vecka 1: Grov planering med främst strukturella konstruktioner. Den interna arkitekturen ska fastslås. Bland annat måste vi bestämma hur back-end och front-end ska kunna kommunicera och samverka harmoniskt. Vi ska lära oss mera om IRC-protokollet och vilka kommandon som är viktiga att implementera.

* Vecka 2: Parser och back-end ska skrivas klart och ha testats någorlunda grundligt. Den grundläggande funktionaliteten ska vara färdig.

* Vecka 3: Front-end ska skrivas klart och användartestning påbörjas och avslutas. Arbetet med vår front-end ska itereras med hjälp av feedback från användartestningen.
