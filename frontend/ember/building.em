class App.BuildingsRoute extends Ember.Route
    model : -> @store.find('building')

class App.BuildingSerializer extends DS.RESTSerializer with CapitalAttrs
    keyForRelation: (key,name) -> key

class App.Building extends DS.Model
    street : DS.attr 'string'
    number : DS.attr 'string'
    city :   DS.attr 'string'
    zip  :   DS.attr 'string'
    lat  :   DS.attr 'string'
    lon  :   DS.attr 'string'
    
    owner_name : DS.attr 'string'
    owner_phone : DS.attr 'string'
    owner_email : DS.attr 'string'

    area : DS.attr 'number'
    
    description : DS.attr 'string'
    
    status : DS.attr 'number'

    comments : DS.hasMany 'comment', async:true

class App.BuildingsController extends Ember.ArrayController
