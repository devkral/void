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
    
    owner_name  : DS.attr 'string'
    owner_phone : DS.attr 'string'
    owner_email : DS.attr 'string'

    area : DS.attr 'number'
    
    description : DS.attr 'string'
    
    status : DS.attr 'number'

    comments : DS.hasMany 'comment', async:true
    newcomment : DS.attr 'string'

class App.BuildingsController extends Ember.ArrayController

class App.BuildingController extends Ember.ObjectController
