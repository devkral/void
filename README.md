# VOID - Vorrichtung zur Organisation von ImmobilienDingen

or in english

# VOID - Collaborative search for immovables

There are some points in life, where you have to organize the search for a new place together
with other people. This might be the search for a new place for your flat-sharing community,
a search for a new office for your company, or a new place for your non-profit. All of those
use cases share, that there are many stakeholders involved in the process of searching,
evaluating, decision-making and so on. Void tries to provide you with a convenient tool to
organize those things as a self hosted solution.

![Screenshot](https://i.imgur.com/2bTfcJ7.png)

## Characteristics

Users of Void cannot register themselves on the site, they must be invited by you. They in turn
can invite other people they consider relevant. Random people cannot simply stop by and
plunder all your carefully collected information.
To invite new users, you can create invitationlinks and give them to your peers.

When you search for immovables, keeping your eyes open is important. But it's even more
efficient if you have other people who keep their eyes open. On Voids startpage, anyone
can enter buildings to suggest them as possible objects of interest for your search.
Simply give the link to your Void-installation to anyone who might be interested in helping
you, and they can notify your group simply by filling out an easy web form.

## Project Status

Unstable. Not featurecomplete yet. Some bugs. In other words, not usable for non-tech-people.
The main work has been written on a spontaneous hackathon on MRMCD13 under the
influence of several tschunks.

## Features

 * [ ] Keep a List of buildings/flats that are relevant to the interests of your group.
 * [ ] Let random people submit buildings which they consider useful for your group.
 * [ ] Invite other users to your group
 * [ ] Have discussions about buildings
 * [ ] Keep track who changed the record of a building/flat when and why.

## License

[AGPLv3](https://www.gnu.org/licenses/agpl.html )

## Building

Get the dependencies together:
 * Working Go-environment for the backend
 * Ember-Script for the frontend
 * Python for the build process
 * MongoDB up and running on the machine you want to deploy to

```shell
git clone https://github.com/grindhold/void
mv void $GOPATH/src/
cd $GOPATH/src/
python2 make.py
```

make.py generates a folder named build, which contains everything you need to deploy on a
webserver.

To start the webserver execute the void binary that you find in the build-folder.
