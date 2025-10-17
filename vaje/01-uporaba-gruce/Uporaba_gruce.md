# Uporaba računske gruče Arnes

Pri predmetu Porazdeljeni sistemi bomo za delo lahko uporabljali računsko gručo [Arnes](https://www.sling.si/arnesova-racunska-gruca/). Trenutno je to drugi najzmogljivejši superračunalnik v Sloveniji ([prvi je Vega](https://si-vegadocs.vega.izum.si/specifikacije/)).
| ![space-1.jpg](slike/arnes.jpg) | 
|:--:| 
| *Računska gruča Arnes*|

### Specifikacije
- 4944 procesorskih jeder
  - 62 x 64 jeder, [AMD Epyc 7702P](https://www.amd.com/en/products/cpu/amd-epyc-7702p)
  - 24 x 12 jeder, [AMD Epyc 7272](https://www.amd.com/en/products/cpu/amd-epyc-7272), 2x [Nvidia V100](https://www.nvidia.com/en-us/data-center/v100/)
  - 7 x 32 jeder, [AMD Epyc 9124](https://www.amd.com/en/products/processors/server/epyc/4th-generation-9004-and-8004-series/amd-epyc-9124.html), 2x [Nvidia H100](https://www.nvidia.com/en-us/data-center/h100/)
  - 3 x 24 jeder, [AMD EPYC 9254](https://www.amd.com/en/products/processors/server/epyc/4th-generation-9004-and-8004-series/amd-epyc-9254.html), 4x [Nvidia H100](https://www.nvidia.com/en-us/data-center/h100/)
- Programska oprema
  - OS [AlmaLinux 9](https://almalinux.org/)
  - Porazdeljen datotečni sistem [Ceph](https://ceph.io/en/)
  - Sistem za upravljanje s posli [SLURM](https://slurm.schedmd.com/)

Do sistema bomo dostopali preko povezave SSH. Uporabniške račune in gesla za dostop najdete na [spletni učilnici FRI](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=37145). Dostop je mogoč samo z uporabo ustreznega ključa SSH in 2FA avtentikacije. Navodila kako si ustvarite lasten ključ in ga dodate v sistem za upravljanje z identitetami najdete [tukaj](https://doc.sling.si/workshops/supercomputing-essentials/02-slurm/06-ssh-key/). Navodila za vzpostavitev 2FA avtentikacije pa najdete [tukaj](https://www.sling.si/dvostopenjska-avtentikacija-za-dostop-do-arnesove-racunske-gruce/).  Po tem, ko ste uredili vse potrebno se na gručo preko ukazne vrstice povežete z ukazom:

```ssh <uporabnisko_ime>@hpc-login.arnes.si```

## Nastavitev okolja

Pri delo z gručo lahko uporabljate poljubno orodje za oddaljen dostop (ukazna lupina, MobaXterm, Putty, FileZilla, WinSCP, CyberDuck, Termius, ...). Priporočamo pa uporabo orodja [VSCode](https://code.visualstudio.com/) v kombinaciji z razširitvijo [Remote - SSH](https://code.visualstudio.com/docs/remote/ssh). Navodila kako vzpostavite povezavo preko VSCode najdete [tukaj](https://doc.sling.si/navodila/vscode/). Pri nastavitvi uporabite Arnesovo vstopno vozlišče `hpc-login.arnes.si`! 

Pri uporabi VSCode so težave s 2FA avtentikacijo. Priporočamo uporabo žetonov krb5 kot je opisano v [navodilih](https://www.sling.si/dvostopenjska-avtentikacija-za-dostop-do-arnesove-racunske-gruce/). Uporabniki operacijskega sistema Windows, lahko žetone generirate v okolju WSL in VSCode naročite, naj za povezovanje uporabi odjemalca ssh znotraj okolja WSL. V poljubni mapi ustvarite skriptno datoteko `ssh.bat` z vsebino:
```
wsl --exec ssh %*
``` 
Nato pa v okolju VSCode spremenite nastavitev `Remote.SSH: Path` tako, da vsebuje absolutno pot do vaše datoteke `ssh.bat`. Vaš privatni ključ in konfiguracijsko datoteko `config`, ki se običajno nahajata v mapi `C:\Users\<uporabnik>\.ssh\` morate prenesti v mapo `~/.ssh` v okolju WSL. Popraviti morate tudi pravice za dostop vašega privatnega ključa:
```
chmod 600 ~/.ssh/<datoteka s privatnim kljucem>
```

## Zaganjanje poslov na gruči

Vodič za delo na gruči in uporamo vmesne programske upreme [SLURM](https://slurm.schedmd.com/) za upravljanje s posli in nalogami najdete [tukaj](https://doc.sling.si/workshops/supercomputing-essentials/01-intro/01-course/). Vsem udeležencem predmeta priporočam, da se prebijejo čez tečaj objavljen na prejšnji povezavi. Pri našem delu z gručo bomo uporabljali rezervacijo `fri`, tako da ne bomo imeli težav s čakanjem, da se naši posli izvedejo. V rezervaciji imamo na voljo nekaj računskih vozlišč, ki jih ostali uporabniki gruče ne morejo zasesti. Rezervacijo pri zaganjanju posla uporabite na naslednji način:

```$ srun --reservation=fri <ime_programa>```

Posle lahko zaganjate tudi s pomočjo opisne skripte bash v kateri navedete zahteve posla. Primer skripte:
```Bash
#!/bin/bash
#SBATCH --job-name=my_job_name
#SBATCH --partition=all
#SBATCH --reservation=fri
#SBATCH --ntasks=4
#SBATCH --nodes=1
#SBATCH --mem-per-cpu=100MB
#SBATCH --output=my_job.out
#SBATCH --time=00:01:00

srun hostname
```
Zgornjo skripto shranite v datoteko končnico `.sh`, npr.: `posel.sh` in jo zaženete z ukazom:
```
$ sbatch posel.sh
```

## Naloga
*Ne šteje kot ena izmed osmih nalog pri predmetu!*
1. Spremeni privzeto uporabniško geslo na [https://fido.sling.si/](https://fido.sling.si/).
2. Ustvari in dodaj ključ SSH v uporabniški profil na [https://fido.sling.si/](https://fido.sling.si/). Navodila najdete na [povezavi](https://doc.sling.si/workshops/supercomputing-essentials/02-slurm/06-ssh-key/).
3. Vzpostavite avtentikacijo s pomočjo krb5 žetonov po [navodilih](https://www.sling.si/dvostopenjska-avtentikacija-za-dostop-do-arnesove-racunske-gruce/).
4. Preko SSH se povežite na vstopno vozlišče Arnes: `hpc-login.arnes.si`.
5. Zaženite program `hostname` na računskem vozlišču znotraj rezervacije `fri`.
6. Zaženite program `nvidia-smi` (izpiše informacije o grafičnih procesnih enotah na vozlišču). Pri zagonu morate uporabiti ustrezno particijo z računskimi vozlišči, ki vsebujejo GPE. (`--partition=gpu`).
7. Kogar zanima malo več gre lahko skozi delavnico [Osnove superračunalništva](https://doc.sling.si/workshops/supercomputing-essentials/01-intro/01-course/).
