class App.SingleMapComponent extends Ember.Component
    lat : ""
    lon : ""
    marker : undefined
    +observer lat,lon
    positionObserver : ->
        if @lat != "" and @lon != ""
            lat = parseFloat @lat
            lon = parseFloat @lon
            if lat == "Nan" or lon == "Nan"
                Ember.Logger.warn "Cant use Latitude: "+@lat+" and Logitude "+@lon
            else
                #if typeof @marker == "object"
                #    Ember.Logger.debug "delete marker"
                #    @map.removeLayer @marker
                @marker = L.marker([lat,lon]).addTo @map
                @map.setView [lat,lon], 15

    didInsertElement : ->
        @_super()
        @$().height "400px"
        @map = L.map this.$().attr 'id'
        #TODO: k-means to find cluster center
        @map.setView [48.3,10.8], 3
        layer = L.tileLayer 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
            attribution: 'Map data &amp; Imagery &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors,     <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
            maxZoom: 18
        layer.addTo(@map)
        @positionObserver()

