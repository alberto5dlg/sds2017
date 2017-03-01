### Subir una rama al repositorio remoto:


    $ git checkout -b nueva-rama
    $ git push -u origin nueva-rama


### Subir cambios de la rama actual:

    (estando en la rama que queremos subir)
    $ git push
El comando git push funcionará correctamente sin más parámetros si previamente hemos subido la rama con un git push -u.


### PullRequest, Merge , y Rebase

El responsable del ticket mezclará el pull request con master desde GitHub. Justo antes de mezclar el pull request, cuando ya se ha tomado la decisión de hacerlo y nadie tiene que subir más cambios, hará un rebase con master para asegurarse que el pull request se introduce en cabeza de master:

    $ git checkout master
    $ git pull
    $ git checkout nueva-rama
    $ git rebase master
    lanzamos los tests para comprobar que todo funciona OK
    y subimos la rama a GitHub (tenemos que usar --force por el rebase)
    $ git push --force


### Descargar una rama del repositorio remoto:

    $ git fetch
    $ git checkout -b nueva-rama origin/nueva-rama
El comando git fetch se descarga todos los cambios pero no los mezcla con las ramas locales. Los deja en ramas cacheadas a las que les da el nombre del servidor y la rama (origin/nueva-rama).

En el caso del comando anterior, una vez cacheada la rama origin/nueva-rama se crea la nueva-rama local con todos sus commits.

### Actualizar una rama con cambios que otros compañeros han subido al repositorio remoto:

      $ git checkout nueva-rama
      $ git pull
El comando git pull es equivalente a un git fetch seguido de un git merge. Algunos recomiendan no usar git pull, sino hacer siempre el merge manual. Por ejemplo:

      $ git checkout nueva-rama
      $ git fetch
      $ git merge origin/nueva-rama
