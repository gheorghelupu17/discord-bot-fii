import discord
import json
from discord.ext import commands
from discord.ext.commands import Bot
TOKEN = ""
#enter guild name or guild id; just one will do :)
GUILD_NAME = 'XXXXX'
GUILD_ID = 1171156323728117780
FILE_NAME = 'memberList'+'.csv'


intents = discord.Intents.all()
client = discord.Client(intents=intents)


@client.event
async def on_ready():
    print(f'{client.user} has connected to Discord!')
    # print(f'{client.guilds} is connected to the following guild:\n')
    for guild in client.guilds:
        print(f'{guild.name}(id: {guild.id==GUILD_ID})\n')
        if guild.id == GUILD_ID:
            print(
                f'{client.user} is scraping the following server: \n' 
                f'{guild.name} (id: {guild.id})'
            )
            roles={}
            for role in guild.roles:
                # line = '{},{},{}\n'.format(member.name+"#"+member.discriminator,member.display_name,member.id)
                roles.update({role.name:f'{role.id}'})
                print(role.name,role.id)        
            jsonRoles = json.dumps(roles)
            with open('roles.json', mode='w',encoding='utf8') as f:
                f.write(jsonRoles)
            with open(FILE_NAME, mode='w',encoding='utf8') as f:
                f.write('username,nickName,id\n')
                for member in guild.members:
                    line = '{},{},{}\n'.format(member.name+"#"+member.discriminator,member.display_name,member.id)
                    f.write(line)


client.run(TOKEN)