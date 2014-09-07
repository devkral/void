class App.InvitationsNewRoute extends Ember.Route
    model : -> Ember.Object.create
                   email: ""
                   link: ""

class App.InvitationRoute extends Ember.Route
    model : (params) -> @store.find 'invitation', params.ident

class App.InvitationSerializer extends DS.RESTSerializer with CapitalAttrs
    keyForRelation: (key,name) -> key

class App.Invitation extends DS.Model
    email : DS.attr 'string'
    name : DS.attr 'string'
    organization : DS.attr 'string'
    password : DS.attr 'string'

class App.InvitationsNewController extends Ember.Controller
    inviteLink : ~> "#{HOST}/#/invitation/#{@content.link}"
    linkPresent : ~> @content.link != ""
    actions:
        invite : ->
            i = @store.createRecord 'invitation',
                  email: @model.email
            self = this
            i.save().then (x) ->
                self.content.link = x.id

class App.InvitationController extends Ember.ObjectController
    actions :
        redeem : ->
            @content.save()
            @transitionTo "index"
