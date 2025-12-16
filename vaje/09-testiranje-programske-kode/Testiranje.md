# Testiranje programske kode v Go

Standardna knjižnica Go vsebuje tudi podporo za izvedbo [testiranja](https://pkg.go.dev/testing) programske kode. Za posamezen paket lahko ustvarimo testne funkcije s katerimi preverimo obnašanje komponent paketa, tako v smislu pravilnosti delovanja kot tudi hitrosti izvajanja.

Znotraj standardne knjižnice Go je na voljo paket `testing`, ki ga uvozimo v datoteki s testno kodo. Datoteki običajno damo ime sestavljeno iz imena datoteke, ki jo testiramo in pripone `_test.go`. Za zgled poskusimo ustvariti nekaj testnih funkcij za [paket](../../predavanja/11-posredovanje-sporocil-2/koda/storage/storage.go) `storage`, ki ste ga spoznali na predavanjih. V mapo `storage` dodamo novo datoteko z imenom `storage_test.go`.

Vanjo dodamo testne funkcije, ki preverjajo funkcionalnost paketa `storage`. Datoteka s testno kodo mora biti del paketa `storage`. Vsaka testna funkcija ima podpis `func TestXxx(t *testing.T)`, kjer je `Xxx` poljubno ime.

Primer testne funkcije, ki preveri delovanje metod `Create` in `Read` paketa storage:
```Go
package storage

import (
	"testing"
)

func TestCreateAndReadAll(t *testing.T) {
	storage := NewTodoStorage()
	todo := Todo{Task: "task1", Completed: false}
    // ustvarimo nov vnos
	if err := storage.Create(&todo, &struct{}{}); err != nil {
		t.Fatalf("create returned error: %v", err)
	}
    // preberemo vse vnose
	got := make(map[string]Todo)
	if err := storage.Read(&Todo{}, &got); err != nil {
		t.Fatalf("read all returned error: %v", err)
	}
    // v shrambi mora biti samo en vnos
	if len(got) != 1 {
		t.Errorf("expected 1 todo, got %d", len(got))
	}
    // preverimo, ce se prebrani vnos ujema z ustvarjenim
	if gotTodo, ok := got[todo.Task]; !ok {
		t.Errorf("expected task %q present", todo.Task)
	} else if gotTodo != todo {
		t.Errorf("read todo mismatch: got %+v want %+v", gotTodo, todo)
	}
}
```
V primeru, da obnašanje testiranih komponent odstopa od pričakovanega sprožimo napako z eno od metod paketa `testing`. V zgornjem primeru sta uporabljeni metodi `Fatalf` in `Errorf`. Prva sproži kritično napako, ki ustavi nadaljnje izvajanje testne funkcije, druga pa sproži napako, vendar se izvajanje testne funkcije nadaljuje.

Test poženemo tako, da se postavimo znotraj mape `storage` in poženemo ukaz:
```
go test
```
Izpiše se status izvajanja testa:
```
PASS
ok      shramba/storage 0.531s
```

Dodajmo še eno testno funkcijo, tokrat za metodo `Update`:
```Go
// pomožna funkcija, za manj ponavljanja znotraj testnih funkcij
func newStorageWithTodos(t *testing.T, todos ...Todo) *TodoStorage {
	t.Helper() // oznaci funkcijo kot pomozno
	storage := NewTodoStorage()
	for _, td := range todos {
		if err := storage.Create(&td, &struct{}{}); err != nil {
			t.Fatalf("create failed for %v: %v", td.Task, err)
		}
	}
	return storage
}
func TestUpdate(t *testing.T) {
	storage := newStorageWithTodos(t, Todo{Task: "present", Completed: false})
    // Testno funkcijo razdelimo na vec podtestov
	t.Run("success", func(t *testing.T) {
		updated := Todo{Task: "present", Completed: true}
		if err := storage.Update(&updated, &struct{}{}); err != nil {
			t.Fatalf("update returned error: %v", err)
		}
		got := make(map[string]Todo)
		if err := storage.Read(&Todo{Task: "present"}, &got); err != nil {
			t.Fatalf("read returned error after update: %v", err)
		}
		if gotTodo := got["present"]; gotTodo.Completed != true {
			t.Errorf("expected updated completion state, got %+v", gotTodo)
		}
	})
	t.Run("not found", func(t *testing.T) {
		err := storage.Update(&Todo{Task: "missing", Completed: true}, &struct{}{})
		if err != ErrorNotFound {
			t.Fatalf("expected error for missing task update, got %v", err)
		}
	})
}
```
Tukaj smo uporabili metodo `Run`, da smo testiranje metode `Update` razbili na dva dela: prvi preverja primer, ko vnos v shrambi obstaja, drugi pa primer, ko vnosa v shrambi ni. Na ta način dobimo bolj strukturiran izpis testiranja.
Sedaj imamo že dve testni funkciji. Da dobimo izpis po posameznih testnih funkcijah in ne samo končnega statusa lahko naš test poženemo z ukazom:
```
go test -v
``` 
in dobimo bolj podroben izpis:
```
=== RUN   TestCreateAndReadAll
--- PASS: TestCreateAndReadAll (0.00s)
=== RUN   TestUpdate
=== RUN   TestUpdate/success
=== RUN   TestUpdate/not_found
--- PASS: TestUpdate (0.00s)
    --- PASS: TestUpdate/success (0.00s)
    --- PASS: TestUpdate/not_found (0.00s)
PASS
ok      shramba/storage 0.460s
```

## Testi Fuzz

Go pozna tudi tako imenovane teste *Fuzz*, ki so predvsem namenjeni lovljenju robnih primerov. Go naključno generira testne podatke glede na začetno seme. Ti testi imajo predpono `Fuzz`. Dodajmo test `Fuzz` za kombinacijo metod `Create`, `Read` in `Delete`:
```Go
func FuzzCreateReadDeleteTest(f *testing.F) {
	// seme za generiranje vhodnih podatkov
	f.Add("seed", false)
	f.Fuzz(func(t *testing.T, task string, completed bool) {
		storage := NewTodoStorage()
		todo := Todo{Task: task, Completed: completed}
		if err := storage.Create(&todo, &struct{}{}); err != nil {
			t.Fatalf("create failed: %v", err)
		}
		got := make(map[string]Todo)
		if err := storage.Read(&Todo{Task: task}, &got); err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("expected one entry, got %d", len(got))
		}
		if got[task].Completed != completed {
			t.Fatalf("completion mismatch: got %v want %v", got[task].Completed, completed)
		}
		if err := storage.Delete(&Todo{Task: task}, &struct{}{}); err != nil {
			t.Fatalf("delete failed: %v", err)
		}
		if err := storage.Read(&Todo{Task: task}, &got); err != ErrorNotFound {
			t.Fatalf("expected not found after delete, got %v", err)
		}
	})
}
```

S pomočjo metode `Add` inicializiramo seme, ki služi Goju kot vodilo pri generiranju testnih podatkov. V našem primeru bomo naključno generirali imena nalog `task` in njihovo stanje `completed`. Metoda `Fuzz` nato naključno generira oba parametra in ju posreduje funkciji, ki smo jo podali kot argument. Test `Fuzz` poženemo s pomočjo ukaza:
```
go test -fuzz .
```
Testiranje se bo izvajalo neskončno časa, oziroma dokler ne najde napake. Čas lahko omejimo s pomočjo stikala `-fuzztime`:
```
go test -fuzz . -fuzztime=5s
```

## Testi zmogljivosti

Včasih nas zanima čas izvajanja posameznih delov kode. V pomoč so nam testi zmogljivosti. Ti imajo predpono Benchmark. Dodajmo test, ki preveri hitrost izvajanja metode `Create`:
```Go
func BenchmarkCreate(b *testing.B) {
	storage := NewTodoStorage()
	b.ResetTimer()
	for b.Loop(){
		if err := storage.Create(&Todo{Task: fmt.Sprintf("task-%d", i)}, &struct{}{}); err != nil {
			b.Fatalf("create failed: %v", err)
		}
	}
}
```
Go samodejno nastavi ustrezno število ponovitev klicev metode `Create`, da zagotovi zanesljivost meritev. Test poženemo s pomočjo ukaza
```
go test -bench .
```
in dobimo izpis:
```
goos: windows
goarch: amd64
pkg: shramba/storage
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkCreate-8        3159963               338.1 ns/op
PASS
ok      shramba/storage 3.447s
```
Iz zgornjega izpisa razberemo, da se je zanka ponovila 3159963x, pri čemer je vsaka iteracija trajala 338,1 ns.

## Domače naloge tokrat ni!
