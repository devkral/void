class App.InvitationsNewRoute extends Ember.Route
    model : -> Ember.Object.create
                   email: ""

class App.InvitationSerializer extends DS.RESTSerializer with CapitalAttrs
    keyForRelation: (key,name) -> key

class App.Invitation extends DS.Model
    email : DS.attr 'string'

class App.InvitationsNewController extends Ember.Controller
    actions:
        invite : ->
            i = @store.createRecord 'invitation',
                  email: @content.email
            i.save()
