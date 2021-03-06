# mcdex - Minecraft Modpack Management

mcdex is a command-line utility that runs on Linux, Windows and OSX 
and makes it easy to manage your modpacks while using the native
Minecraft launcher.

## Downloads

You can find the most recent releases here:

* [Windows](http://files.mcdex.net/releases/win32/mcdex.exe)
* [Linux](http://files.mcdex.net/releases/linux/mcdex)
* [OSX](http://files.mcdex.net/releases/osx/mcdex)

## Getting started

First, make sure you have the most recent database of mods:

```
mcdex db.update
```

## Listing available mod packs on Curseforge

If you want to find all the published modpacks available with 'engineer' in the name, you can do:

```
mcdex pack.list engineer 
```

You'll see output like:

```
advanced-engineering-immersive-rocketry | A RotaryCraft-esque pack for 1.12 | 154 downloads
age-of-engineering | Tech-only modpack with guide and trophies. Focus in progression, divided in 15 ages | 438,304 downloads
ancient-engineering | A light pack focused around Ancient Warfare and Immersive Engineering. | 4,856 downloads
arie-advanced-engineering | Modpack filled with immersive tech mods, combined with the possibility of exploring the galaxy! | 53 downloads
atmospheric-engineering | an overloaded skyfactory style modpack | 0 downloads
coles-engineering | An engineering modpack with big dig vibes | 60 downloads
...
...
```

The first part of each line is the "slug"; we'll use this in the next step to install the modpack.

## Installing a modpack from Curseforge

Now, let's install the [Age of Engineering](https://minecraft.curseforge.com/projects/age-of-engineering) modpack.

Using the slug we found via pack.list, run:
```
mcdex pack.install aoe age-of-engineering
```

Note that we provide the name "aoe"; mcdex uses this to install the modpack into your Minecraft home directory
under ```<minecraft>/mcdex/pack/aoe```. Alternatively, you can control what directory the modpack is installed in by passing
an absolute path (e.g. `c:\aoe` or `/Users/dizzyd/aoe`) as the name and mcdex will use that instead.

Once the install is done, you can fire up the Minecraft launcher and you should have a new profile for the aoe pack!

## Creating a new modpack

We can start a new modpack by using the ```pack.create``` command:

```
mcdex pack.create mypack 1.11.2
```

Note that the recommended version of Forge is installed automatically. If you want to force a specific Forge to be used,
you can do
```
mcdex pack.create mypack 1.11.2 13.20.1.2386
```

As before, since we passed a non-absolute filename - 'mypack' - the pack will be created under the Minecraft home directory. 
If we wanted to create the modpack in our home directory (on OSX), we would do:

```
mcdex pack.create /Users/dizzyd/mypack 1.11.2
```

```pack.create``` will create the directory, make sure the appropriate version of Forge is installed and start a manifest.json. 
In addition, it will create an entry in the Minecraft launcher so you can launch the pack.

## Installing individual mods

Once you have a modpack, either installed from CurseForge or one you created locally, you can add individual mods to it. Let's
add [Immersive Engineering](https://minecraft.curseforge.com/projects/immersive-engineering) to our new pack:

```
mcdex mod.select mypack immersive-engineering
```

This will search the database of mods for one named 'immersive-engineering' and find the most recent stable version and
add it to the pack's manifest.jason. To actually install the mod, you need to install the pack:

```
mcdex pack.install mypack
```

## Listing available mods

If you want to find all the mods with 'Map' in the name, you can do:

```
mcdex mod.list Map
```

Alternatively, if you want to only look for mods with 'Map' that work on 1.10.2, you can do:

```
mcdex mod.list Map 1.10.2
```

## Updating mods within a pack

If you want to update all the mods within a pack, you can now run:

```
mcdex mod.update.all mypack
```

This will walk over all the installed mods within a manifest and look for more recent versions of the mod that work with your selected version of minecraft.
It will list all the affected mods as it runs, so you can see exactly what changed. If you want to make sure a mod doesn't get updated, in the manifest.json
you can add a "locked": true entry to that mod and mcdex will not upgrade it.

Note that you can also run this command with a -n flag (dry run) so that it will simply print out the mods that were 
updated without actually updating the manifest.

Once you've updated the manifest with mod.update.all, you need to re-install the pack to make sure the new mods are updated:

```
mcdex pack.install mypack
```
