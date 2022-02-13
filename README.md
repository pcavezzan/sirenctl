sirenctl ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/pcavezzan/sirenctl/build)
========

Petit utilitaire permettant de récupérer l'établissement (SIRET) associée à l'entreprise (SIREN) et le code postal de l'établissement ciblé.

## Usage


### A l'aide un SIREN et code postal

```shell
./sirenctl --siren=791092174 --codepostal=95200         
siren;code_postal;siret
791092174;95200;79109217400014
```

### A l'aide d'un fichier

```shell
./sirenctl -f ./example/input.csv -o ./output/result.csv
```
Un fichier `result.csv` est alors créé dans le répertoire `output`:
```shell
cat output/result.csv 
siren;code_postal;siret
791092174;95200;79109217400014
443061841;75009;44306184100047
443061841;75009;44306184100062
443061841;75009;44306184100070
443061841;75009;44306184100088
```

