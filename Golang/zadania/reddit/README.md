# Fetcher

Celem zadania jest napisanie aplikacji, która pobierze listę wpisów z portalu reddit, i umożliwi zapisanie do pliku, bufora lub innych structów spełniających interfejs `io.Writer`

* Adres URL do odpytywania API:
https://www.reddit.com/r/golang.json
* Struct i interfejs znajdują się w pliku `redditFetcher.go`

### Wymagania:
* Aplikacja musi poprawnie pobierać, (de)serializować i zapisywać dane z redditowego API:
  * wyeksportowane pole Title o typie string
  * wyeksportowane pole URL o typie string
* Struct musi mieć również implementować interfejs `RedditFetcher`
* Podanego interfejsu nie można zmieniać, ale można rozszerzać strukturę/interfejs (dla chętnych).
* Obsłużyć błędy np. poprzez logowanie

#### W celu zaliczenia zadania należy:
1. Zaimplementować własny struct i metody tak by spełniać interfejs RedditFetcher
2. Zapisać pobrany listing do wybranego przez siebie io.Reader'a (plik lub bufor)


#### Dodatkowo można poeksperymentować:
1. Klient używa Context z timeoutem
2. Współbieżne pobieranie innych wybranych subredditów SFW