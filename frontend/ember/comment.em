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

class App.CommentSerializer extends DS.RESTSerializer with CapitalAttrs

class App.Comment extends DS.Model
    text: DS.attr 'string'
    date : DS.attr 'string'
    user : DS.belongsTo 'user', async:true
    type : DS.attr 'string'
    building: DS.belongsTo 'building', async:true

