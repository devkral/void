class App.CommentSerializer extends DS.RESTSerializer with CapitalAttrs

class App.Comment extends DS.Model
    text: DS.attr 'string'
    date : DS.attr 'string'
    user : DS.belongsTo 'user', async:true
    type : DS.attr 'string'
    building: DS.belongsTo 'building', async:true

