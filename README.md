# Projektplan
Hej
## Programbeskrivning
Programmet ska kunna koppla upp sig på populära IRC-servrar och kunna stödja grundläggande operationer i IRC-protokollet. Med hjälp av Go:s multitrådning hoppas vi också kunna implementera ett smidigt sätt att samtidigt vara uppkopplad mot flera kanaler och/eller servrar.

## Användarbeskrivning
Vår målgrupp är ganska vana datoranvändare som redan har viss erfarenhet av att använda andra IRC-klienter. Vi kommer dock att sikta på att implementera ett hjälpkommando som kan täcka åtminstone den grundläggande funktionaliteten.

## Användarscenarion

### I
Användaren går in i vårt program och ser två textrutor och en send-knapp. I den översta rutan står det ”Welcome to irken. For help type /help”. Vår användar skriver ”/help” i den nedre textrutan och trycker sedan på ”send”. Han får då upp en lista på kommandon han kan använda. Han väljer att gå in på QuakeNet och chatta på kanalen #kthd. Han skriver därför ”/connect quakenet” och sedan ”/join #kthd”. För att folk ska känna igen honom skriver han ”/nick crazybanana92” och blir då ett med sitt alterego. Han skriver sitt första inlägg ”vim drools, emacs RULES!”, trycker på ”send” och blir sedermera utkastad ur kanalen.

### II
Användaren går tillbaka in i #kthd och tar sig samman. Han vill nu chatta i en annan kanal, men han vill samtidigt ta emot meddelanden från #kthd. Han går upp i en meny och väljer ”ny tabb”. Han hamnar då i ett nytt fönster, men har kvar sin uppkoppling mot QuakeNet och sitt fina nick. Han väljer att gå in på en annan kanal. Han skriver då ”/join #fabhack” och pratar en stund på den nya kanalen. Efter tag vill han tillbaka till sin förra kanal. Han väljer då ”#kthd”-tabben från en meny, och kommer tillbaka till den förra kanalen. Han kan då se vad alla andra har skrivit medan han var borta.

## Testplan
Självklart kommer programmet att testas med både manuella och automatiska tester av oss som programmerar. Vi kommer försöka att skriva så mycket testkod vi kan, men vet inte hur lätt detta blir eftersom stor del av programmet består av nätverks- och grafikkod. Testanvändaren kommer att genomföra våra användarscenarion och ge feedback på hur åtkomliga dessa var att göra.

