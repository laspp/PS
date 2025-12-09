# Uporaba gRPC v programskem jeziku go

[gRPC](https://grpc.io/) je Googlova rešitev za prenosljivo, odprtokodno ogrodje za izvedbo klicev oddaljenih metod. Za prenos podatkov uporablja protokol HTTP/2, pri čemer podatke kodira z  binarnim protokolom [Protocol Buffers](https://protobuf.dev/).

Na [primeru kode](../../predavanja/11-posredovanje-sporocil-2/koda/grpc/), ki ste jo obravnavali na predavanjih, si poglejmo kako dodati novo metodo v strežnik CRUD. V strežnik bomo dodali metodo `ReadAll`, ki bo izpisala vse vnose v shrambi.

Uporabili bomo pristop s podatkovnimi tokovi. gRPC namreč poleg izvajanja posameznih zahtevkov ponuja tudi možnost komunikacije preko tokov. Ti so uporabni, kadar želimo med odjemalcem in strežnikom prenesti večje število sporočil.

Najprej si namestimo prevajalnik za `Protocol Buffers` protoc. Navodila za različne OS najdete [tukaj](https://protobuf.dev/installation/). Nato namestimo še gRPC po [navodilih](https://grpc.io/docs/languages/go/quickstart/).
Če smo na gruči Arnes sta Go in Protocol Buffers že nameščena; najprej naložimo modula `Go` in `protobuf` ter ustrezno verzijo `binutils`:

```Bash
$ module load Go
$ module load protobuf/23.0-GCCcore-12.2.0
$ module load binutils/2.39-GCCcore-12.2.0
```

Na gruči Arnes se uporablja nekoliko starejša verzija Go in Protocol Buffer, zato namestimo starejšo verzijo gRPC. Če do sedaj še nismo, moramo (samo prvič) namestiti paketa `protobuf` in `grpc` za **go**:

```Bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

Nato ustrezno nastavimo pot do prevajalnika za **go**:
```
$ export PATH="$PATH:$(go env GOPATH)/bin"
```

Posodobimo opisno datoteko `protobufStorage/protobufStorage.proto`, kjer so definirane storitve našega strežnika. V datoteko dodamo novo metodo `ReadAll`:
```protobuf
service CRUD {
    rpc Create(Todo) returns (google.protobuf.Empty);
    rpc Read(Todo) returns (TodoStorage);
    rpc Update(Todo) returns (google.protobuf.Empty);
    rpc Delete(Todo) returns (google.protobuf.Empty);
    // nova metoda ReadAll
    rpc ReadAll(google.protobuf.Empty) returns (stream Todo); 
}
```

Ker nova metoda ne potrebuje argumentov, vendar jih gRPC zahteva, uporabimo vgrajeni tip `google.protobuf.Empty`, ki predstavlja prazno sporočilo. Metoda vrača tok sporočil (stream) tipa Todo.

Sedaj uporabimo prevajalnik za pretvorbo opisne datoteke `.proto` v kodo, ki jo uvozimo v go. V ukazni vrstici se prestavimo v mapo `protobufStorage` in prevedemo posodobljeno opisno datoteko:
```Bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protobufStorage.proto
```

Prevajalnik ustvari dve novi datoteki `protobufStorage.pb.go` in `protobufStorage_grpc.pb.go`. Prva vsebuje vse potrebne metode in podatkovne tipe za kodiranje in dekodiranje sporočil. Druga pa definira vmesnik RPC strežnika in odjemalca.

## Strežnik

Najprej dopolnimo kodo strežnika, ki se nahaja v datoteki `grpc/streznik.go` in dodajmo implementacijo naše nove metode `ReadAll`.
V datoteki `protobufStorage_grpc.pb.go` najdemo generiran prototip funkcije:

```Go
ReadAll(*emptypb.Empty, grpc.ServerStreamingServer[Todo]) error
```
Metoda vzame dva argumenta, prvi je prazno sporočilo (tega smo definirali pri opisu metode), drugi pa je vmesnik za podatkovni tok, ki vključuje metodo `Send`, s pomočjo katere bomo pošiljali podatke odjemalcu.

Na koncu datoteke `streznik.go` pripravimo telo funkcije `ReadAll`:
```Go
func (s *serverCRUD) ReadAll(e *emptypb.Empty, stream grpc.ServerStreamingServer[protobufStorage.Todo]) error {
	dict := make(map[string](storage.Todo))
	err := s.todoStore.Read(&storage.Todo{}, &dict)
	for _, todo := range dict {
		todo_temp := &protobufStorage.Todo{Task: todo.Task, Completed: todo.Completed}
		if err := stream.Send(todo_temp); err != nil {
			return err
		}
	}
	return err
}
```
V funkciji najprej ustvarimo prazen slovar, ki bo hranil vse vnose, ki jih želimo poslati odjemalcu. Nato preberemo vse vnose iz pomnilniške shrambe ter jih zaporedoma pošljemo odjemalcu.

## Odjemalec

Dopolnimo še kodo odjemalca v datoteki `odjemalec.go`. Najprej dodajmo še en vnos v našo shrambo:
```Go
lecturesCreate := protobufStorage.Todo{Task: "predavanja", Completed: false}
lecturesUpdate := protobufStorage.Todo{Task: "predavanja", Completed: true}
homeworkCreate := protobufStorage.Todo{Task: "domaca naloga", Completed: true} // dodamo domačo nalogo
practicals := protobufStorage.Todo{Task: "vaje", Completed: false}
readAll := protobufStorage.Todo{Task: "", Completed: false}
```

Odstranimo naslednji odsek kode, ki izpiše vsebino shrambe na strežniku in jo nadomestimo z dodajanjem novega vnosa v bazo:
```Go
// preberemo vse zapise
//	fmt.Print("7. Read *: ")
//	if response, err := grpcClient.Read(contextCRUD, &readAll); err == nil {
//		fmt.Println(response.Todos, ": done")
//	} else {
//		panic(err)
//	}

// ustvarimo zapis
fmt.Print("7. Create: ")
if _, err := grpcClient.Create(contextCRUD, &homeworkCreate); err != nil {
    panic(err)
}
fmt.Println("done")
```

Sedaj uporabimo našo novo oddaljeno metodo in izpišimo vse vnose v shrambi na strežniku. Najprej poiščemo prototip metode v datoteki `protobufStorage_grpc.pb.go`:
```Go
ReadAll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Todo], error)
```
Metoda za argumete vzame trenutni kontekst, prazno sporočilo ter dodatne opcijske argumente, ki pa jih v našem primeru nimamo. Metoda vrne vmesnik `grpc.ServerStreamingClient[Todo]`, ki ima naslednjo definicijo:
```Go
type ServerStreamingClient[Res any] interface {
    Recv() (*Res, error)
    ClientStream
}
```
Vmesnik vsebuje metodo `Recv`, ki služi prejemanju sporočil po podatkovnem toku. Naslednji odsek kode prebere vse vnose v shrambi strežnika in jih izpiše:

```Go
// preberemo vse zapise stream
fmt.Print("8. Read *: \n")
if stream, err := grpcClient.ReadAll(contextCRUD, &emptypb.Empty{}); err == nil {
    for {
        todo, err := stream.Recv()
        if err == io.EOF {
            break
        }
        fmt.Println(todo)
    }
} else {
    panic(err)
}
fmt.Println("done")
```
V stavek import moramo dodati še paketa `io` in `"google.golang.org/protobuf/types/known/emptypb"`, ki ju uporabljamo znotraj zgornjega odseka kode. Celotno kodo najdete [tukaj](./koda/resitev/).

## Preizkus

Sedaj lahko poženemo strežnik. Preverimo, da se nahajamo v mapi `grpc` in izvedemo ukaz, da zaženemo strežnik:
```Bash
$ go run .
```
 in odjemalca:

 ```Bash
$ go run . -s localhost -p 9876
```
Pri zagonu odjemalca moramo za ime strežnika in vrata vpisati vrednosti, ki jih izpiše strežnik ob zagonu. V primeru, da poganjamo strežnik in odjemalca na gruči Arnes, oba ukaza izvedemo preko SLURM z ukazom `srun`.

# Domača naloga 7

Vaša naloga je posodobiti storitev CRUD iz zgornjega primera tako, da bo nudila naročnino na dogodke v shrambi "TodoStorage". Rešitev (modul) v obliki datoteke zip oddajte preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=60917).

Odjemalcu želimo omogočiti spremljanje vseh dogodkov, ki se zgodijo v naši storitvi CRUD. Zato bomo specifikacijo programskega vmesnika razširili z naslednjimi vnosi:
```protobuf
message TodoEvent {
    Todo    T = 1;
    string  Action = 2;
}

service CRUD {
    //...
    rpc Subscribe(google.protobuf.Empty) returns (stream TodoEvent);
}
```
Sporočilo `TodoEvent` vsebuje vnos `T`, ki se je v shrambi spremenil in tip izvedene akcije `Action` v obliki niza:
- `create`: Dodan nov vnos v shrambo.
- `update`: Posodobljen vnos v shrambi.
- `delete`: Izbrisan vnos.

Odjemalec se s pomočjo metode `Subscribe` naroči na dogodke in jih prejme v obliki sporočila `TodoEvent` iz strežnika takoj, ko se zgodijo. Delovanje naročnine preverite tako, da znotraj odjemalca v ločeni gorutini zaženete več akcij, ki spreminjajo vsebino shrambe. Med posamezne klice dajte kratko pavzo, (1 sekunda). V drugi gorutini pa spremljate zaporedje dogodkov s pomočjo metode Subscribe in jih izpisujte.

**Rok za oddajo 16. 12. 2025.**

