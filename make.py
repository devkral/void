#!/usr/bin/python

########################################################################
#
#    © 2014 Daniel 'grindhold' Brendle and Contributors
# 
#    This file is part of Void.
# 
#    Void is free software: you can redistribute it and/or
#    modify it under the terms of the GNU Affero General Public License
#    as published by the Free Software Foundation, either
#    version 3 of the License, or (at your option) any later
#    version.
# 
#    Void is distributed in the hope that it will be
#    useful, but WITHOUT ANY WARRANTY; without even the implied
#    warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
#    PURPOSE. See the GNU Affero General Public License for more details.
# 
#    You should have received a copy of the GNU Affero General Public
#    License along with Void.
#    If not, see http://www.gnu.org/licenses/.
#
#########################################################################


import shutil
import os
import sys
import subprocess

def compile_go(debug=True):
    output = subprocess.check_output("""cd backend ; 
                                        go get ; 
                                        go build -o void ;
                                        cd ..
                                        exit 0""", stderr = subprocess.STDOUT, shell=True)
    print output
    return

def buildframe():
    readframe = open('frontend/application.html', 'r')
    writeframe = open('build/index.html', 'w')
    for line in readframe:
        if line.find('<body>') == -1:
            writeframe.write(line)
        else:
            writeframe.write(line)
            insertcomponents(writeframe)
            inserttemplates(writeframe)

def copythings():
    shutil.copytree('frontend/dist', 'build/static/', symlinks=False, ignore=None)
    shutil.copytree('frontend/img', 'build/static/img', symlinks=False, ignore=None)
    shutil.copy2('backend/void', 'build')

def inserttemplates(writeframe):
    templatenames = os.listdir('frontend/templates/')
    for i in templatenames:
        if not i.endswith(".hbs"):
            continue
        if i != 'application.hbs':
            writeframe.write(' '*4 + '<script type="text/x-handlebars" data-template-name="' \
                              + i.split('.')[0].replace('-', '/') + "\">\n")
        else:
            writeframe.write(' '*4 + '<script type="text/x-handlebars">\n')
        frame = open('frontend/templates/' + i)
        for line in frame:
            writeframe.write(' '*8 + line)
        writeframe.write(' '*4 + '</script>\n')

def insertcomponents(writeframe):
    componentnames = os.listdir('frontend/components/')
    for i in componentnames:
        if not i.endswith(".hb"):
            continue
        writeframe.write(' '*4 + '<script type="text/x-handlebars" data-template-name="components/' + \
                          i.split('.')[0]+"\">\n")
        frame = open('frontend/components/' + i)
        for line in frame:
            writeframe.write(' '*8 + line)
        writeframe.write(' '*4 + '</script>\n')

def packproject():
    if os.path.exists('build'):
        shutil.rmtree('build')
    os.mkdir('./build')
    buildframe()
    copythings()

def compile_emberscript(debug=True):
    code = ""
    em_folder = "frontend/ember"
    comp_filename = em_folder+"/__drawn_together.em"
    try:
      os.unlink(comp_filename)
    except: pass
    for embersource in os.listdir(em_folder):
      if embersource.endswith(".em"):
        f = open(em_folder+"/"+embersource)
        if embersource == "application.em":
            code = f.read()+code
        else:   
            code += f.read()
        f.close()
    comp_file = open(comp_filename,"w")
    comp_file.write(code)
    comp_file.close()
    
    compiler_arguments = ['ember-script','-j','-i '+comp_filename]
    if not debug:
      compiler_arguments.append('-m')

    x = subprocess.check_output("ember-script -j -i %s %s; exit 0"%( \
            comp_filename, ('-m','')[debug]), stderr = subprocess.STDOUT, shell=True)
    out = open("frontend/dist/js/application-0.1.js","w")
    out.write(x)
    out.close()

if __name__ == '__main__':
    debug = True
    if len(sys.argv) > 1 and sys.argv[1] == "--production":
        debug = False
    compile_go(debug)
    compile_emberscript(debug)
    packproject()
