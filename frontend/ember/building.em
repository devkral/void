class App.BuildingsRoute extends Ember.Route
    model : -> @store.find 'building'

class App.BuildingsMapRoute extends Ember.Route
    
class App.BuildingsListRoute extends Ember.Route

class App.BuildingRoute extends Ember.Route
    model : (params) -> @store.find 'building', params.ident

class App.BuildingViewRoute extends Ember.Route

class App.BuildingSerializer extends DS.RESTSerializer with CapitalAttrs
    keyForRelation: (key,name) -> key

class App.Building extends DS.Model
    street : DS.attr 'string'
    number : DS.attr 'string'
    city   : DS.attr 'string'
    zip    : DS.attr 'string'
    lat    : DS.attr 'string'
    lat_f  : DS.attr 'number'
    lon    : DS.attr 'string'
    lon_f  : DS.attr 'number'
    
    ownername  : DS.attr 'string'
    ownerphone : DS.attr 'string'
    owneremail : DS.attr 'string'

    area : DS.attr 'number'
    
    description : DS.attr 'string'
    
    status : DS.attr 'number'

    comments : DS.hasMany 'comment', async:true
    newcomment : DS.attr 'string'

class App.BuildingsController extends Ember.ArrayController

class App.BuildingController extends Ember.ObjectController
    oldData : null
    editMode : ~> false
    actions:
        post: ->
            self = this
            @store.find('user', name: App.sessionUser()).then (u) ->
                u = u.objectAt 0
                c = self.store.createRecord 'comment',
                      text: self.content.newcomment
                      user: u
                      building: self.content
                self.content.comments.addObject c
                c.save()
        enterEdit : ->
            @oldData = Ember.Object.create
                street : @content.street
                number : @content.number
                city   : @content.city
                zip    : @content.zip
                ownername  : @content.ownername
                ownerphone : @content.ownerphone
                owneremail : @content.owneremail
                area : @content.area
                description : @content.description
                status : @content.status
            @editMode = true
            return
        leaveEdit : ->
            @content.street = @oldData.street
            @content.number = @oldData.number
            @content.city = @oldData.city
            @content.zip = @oldData.zip
            @content.ownername = @oldData.ownername
            @content.ownerphone = @oldData.ownerphone
            @content.owneremail = @oldData.owneremail
            @content.area = @oldData.area
            @content.description = @oldData.description
            @content.status = @oldData.status
            @editMode = false
            return
        executeEdit : ->
            @content.save()
