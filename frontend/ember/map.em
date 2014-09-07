###########################################################
# © 2014 Daniel 'grindhold' Brendle and Contributors 
#
# This file is part of Void.
#
# Void is free software: you can redistribute it and/or
# modify it under the terms of the GNU Affero General Public License
# as published by the Free Software Foundation, either
# version 3 of the License, or (at your option) any later
# version.
#
# Void is distributed in the hope that it will be
# useful, but WITHOUT ANY WARRANTY; without even the implied
# warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
# PURPOSE. See the GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public
# License along with Void.
# If not, see http://www.gnu.org/licenses/.
###########################################################

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
        @$().height "150px"
        @map = L.map this.$().attr 'id'
        #TODO: k-means to find cluster center
        @map.setView [48.3,10.8], 3
        layer = L.tileLayer 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
            attribution: 'Map data &amp; Imagery &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors,     <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
            maxZoom: 18
        layer.addTo(@map)
        @positionObserver()

class App.MultiMapComponent extends Ember.Component
    buildings : []
    markers : []
    +observer buildings.[]
    positionObserver : ->
        self = this
        @buildings.forEach (building) ->
          if building.lat != "" and building.lon != ""
              m = L.marker([building.lat_f,building.lon_f]).addTo self.map
              m.bindPopup """<a href="/#/building/#{building.id}/view">#{building.street}</a>"""
              self.markers.push m
              #TODO: k-means to find cluster center

    didInsertElement : ->
        @_super()
        @$().height "600px"
        @map = L.map this.$().attr 'id'
        #TODO: k-means to find cluster center
        @map.setView [48.3,10.8], 3
        layer = L.tileLayer 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
            attribution: 'Map data &amp; Imagery &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors,     <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>'
            maxZoom: 18
        layer.addTo(@map)
        @positionObserver()

