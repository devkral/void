class App.UserSerializer extends DS.RESTSerializer with CapitalAttrs
    keyForRelation: (key,name) -> key

class App.User extends DS.Model
    name : DS.attr 'string'
    email : DS.attr 'string'
    organization : DS.attr 'string'
