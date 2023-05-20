(()=>{"use strict";var t,e={312:(t,e)=>{var n;Object.defineProperty(e,"__esModule",{value:!0}),e.EVENTS_NAME=void 0,(n=e.EVENTS_NAME||(e.EVENTS_NAME={})).chestLoot="chest-loot",n.attack="attack"},325:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Actor=void 0;var s=function(t){function e(e,n,o,i,s){var r=t.call(this,e,n,o,i,s)||this;return r.hp=100,r.speed=200,r.speedUp=100,e.add.existing(r),e.physics.add.existing(r),r.getBody().setCollideWorldBounds(!0),r}return i(e,t),e.prototype.getDamage=function(t){t&&(this.hp=this.hp-t)},e.prototype.getHPValue=function(){return this.hp},e.prototype.checkFlip=function(){this.body.velocity.x<0?this.scaleX=-1:this.scaleX=1},e.prototype.getBody=function(){return this.body},e}(n(260).Physics.Arcade.Sprite);e.Actor=s},148:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Chest=void 0;var s=n(260),r=n(312),a=function(t){function e(e,n,o,i){var s=t.call(this,e,n,o,"tiles_spr",595)||this;return s.setScale(1.5),e.add.existing(s),e.physics.add.existing(s),s.chestId=i,s.initCollision(),s}return i(e,t),e.prototype.initCollision=function(){var t=this,e=this.scene.player;this.scene.physics.add.overlap(e,this,(function(e,n){t.scene.game.events.emit(r.EVENTS_NAME.chestLoot),n.destroy(),window.connection.send({id:t.chestId,messageType:"chest_grab"})}))},e.prototype.deleteChest=function(){this.destroy(!0)},e}(s.Physics.Arcade.Sprite);e.Chest=a},321:function(t,e,n){var o=this&&this.__assign||function(){return o=Object.assign||function(t){for(var e,n=1,o=arguments.length;n<o;n++)for(var i in e=arguments[n])Object.prototype.hasOwnProperty.call(e,i)&&(t[i]=e[i]);return t},o.apply(this,arguments)};Object.defineProperty(e,"__esModule",{value:!0}),e.Connection=e.uid=void 0;var i=n(815),s=n(753),r=n(148);e.uid=function(){return Date.now().toString(36)+Math.random().toString(36).substr(2)};var a=function(){function t(){var t=this;this.syncObjects=new Map,this.localStorage=window.sessionStorage,this.initLocalData(),setTimeout((function(){t.initConnection(),t.initCallbacks()}),1500)}return t.prototype.send=function(t){var e;if((null===(e=this.socket)||void 0===e?void 0:e.readyState)===WebSocket.OPEN){var n={messageType:t.messageType,data:JSON.stringify(t)};this.socket.send(JSON.stringify(n))}},t.prototype.close=function(){this.socket.readyState===WebSocket.OPEN&&this.socket.close()},t.prototype.initScene=function(t){return this.scene=t,new s.Player(t,this.gamePlayerData.coordinates.x,this.gamePlayerData.coordinates.y,this.gamePlayerData.id)},t.prototype.initLocalData=function(){var t=this.localStorage.getItem("game_player_data");t?this.gamePlayerData=JSON.parse(t):(this.gamePlayerData={id:(0,e.uid)(),coordinates:{x:200,y:600,directionX:1},health:100,directionX:1},this.localStorage.setItem("game_player_data",JSON.stringify(this.gamePlayerData)))},t.prototype.initConnection=function(){var t=window.location.hostname,e="",n="wss";"localhost"===t?(e=":8080",n="ws"):window.location.hostname.includes("github.io")&&(t="tanks-maciejdominiak.b4a.run");var o="".concat(n,"://").concat(t).concat(e,"/tanks/objects:exchange");this.socket=new WebSocket(o)},t.prototype.initCallbacks=function(){var t=this;this.socket.onopen=function(e){t.socket.send(JSON.stringify(o(o({},t.gamePlayerData),{messageType:"create_player"}))),console.log("connected")},this.socket.onmessage=function(e){try{var n=JSON.parse(e.data);switch(console.log(n),n.messageType){case"create_player":t.syncObjects.set(n.id,new i.OtherPlayer(t.scene,n.coordinates.x,n.coordinates.y,n.id));break;case"user_disconnected":t.syncObjects.has(n.id)&&(t.syncObjects.get(n.id).deletePlayer(),t.syncObjects.delete(n.id));break;case"status":t.syncObjects.has(n.id)?t.syncObjects.get(n.id)instanceof i.OtherPlayer&&t.syncObjects.get(n.id).updatePlayer(n.coordinates.x,n.coordinates.y,n.coordinates.directionX,n.health):t.syncObjects.set(n.id,new i.OtherPlayer(t.scene,n.coordinates.x,n.coordinates.y,n.id));break;case"user_attack":t.syncObjects.get(n.id)instanceof i.OtherPlayer&&t.syncObjects.get(n.id).attack();break;case"create_chest":t.syncObjects.set(n.id,new r.Chest(t.scene,n.coordinates.x,n.coordinates.y,n.id));break;case"chest_destroy":t.syncObjects.get(n.id)instanceof r.Chest&&(t.syncObjects.get(n.id).deleteChest(),t.syncObjects.delete(n.id))}}catch(t){console.log(t)}},this.socket.onerror=function(t){console.log(t)},this.socket.onclose=function(t){console.log("closed")}},t}();e.Connection=a},97:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Enemy=void 0;var s=n(325),r=n(312),a=function(t){function e(e,n,o,i,s,a){var c=t.call(this,e,n,o,i,a)||this;return c.AGRESSOR_RADIUS=100,c.target=s,e.add.existing(c),e.physics.add.existing(c),c.getBody().setSize(16,16),c.getBody().setOffset(0,0),c.attackHandler=function(){Phaser.Math.Distance.BetweenPoints({x:c.x,y:c.y},{x:c.target.x,y:c.target.y})<c.target.width&&(c.getDamage(),c.disableBody(!0,!1),c.scene.time.delayedCall(300,(function(){c.destroy()})))},c.scene.game.events.on(r.EVENTS_NAME.attack,c.attackHandler,c),c.on("destroy",(function(){c.scene.game.events.removeListener(r.EVENTS_NAME.attack,c.attackHandler)})),c}return i(e,t),e.prototype.preUpdate=function(){Phaser.Math.Distance.BetweenPoints({x:this.x,y:this.y},{x:this.target.x,y:this.target.y})<this.AGRESSOR_RADIUS?(this.getBody().setVelocityX(this.target.x-this.x),this.getBody().setVelocityY(this.target.y-this.y)):this.getBody().setVelocity(0)},e.prototype.setTarget=function(t){this.target=t},e}(s.Actor);e.Enemy=a},815:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.OtherPlayer=void 0;var s=n(325),r=n(67),a=function(t){function e(e,n,o,i){var s=t.call(this,e,n,o,"king")||this;return s.getBody().setSize(30,30),s.getBody().setOffset(8,0),s.hpValue=new r.Text(s.scene,s.x,s.y-s.height,s.hp.toString()).setFontSize(12).setOrigin(.8,.5),s.playerId=i,s}return i(e,t),e.prototype.updatePlayer=function(t,e,n,o){this.setPosition(t,e),this.scaleX=n,this.hp=o,this.hpValue.setPosition(this.x,this.y-.4*this.height),this.hpValue.setOrigin(.8,.5)},e.prototype.attack=function(){this.anims.play("attack",!0)},e.prototype.deletePlayer=function(){this.hpValue.destroy(!0),this.destroy(!0)},e.prototype.getDamage=function(e){t.prototype.getDamage.call(this,e),this.hpValue.setText(this.hp.toString())},e}(s.Actor);e.OtherPlayer=a},753:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Player=void 0;var s=n(325),r=n(67),a=n(312),c=function(t){function e(e,n,o,i){var s=t.call(this,e,n,o,"king")||this;return s.keyW=s.scene.input.keyboard.addKey("W"),s.keyA=s.scene.input.keyboard.addKey("A"),s.keyS=s.scene.input.keyboard.addKey("S"),s.keyD=s.scene.input.keyboard.addKey("D"),s.keyShift=s.scene.input.keyboard.addKey("Shift"),s.keySpace=s.scene.input.keyboard.addKey(32),s.keySpace.on("down",(function(t){s.anims.play("attack",!0),s.scene.game.events.emit(a.EVENTS_NAME.attack),s.sendAttack()})),s.getBody().setSize(30,30),s.getBody().setOffset(8,0),s.hpValue=new r.Text(s.scene,s.x,s.y-s.height,s.hp.toString()).setFontSize(12).setOrigin(.8,.5),s.playerId=i,s.initAnimations(),s}return i(e,t),e.prototype.update=function(){var t,e,n,o,i;this.getBody().setVelocity(0);var s=(null===(t=this.keyShift)||void 0===t?void 0:t.isDown)?this.speed+this.speedUp:this.speed;(null===(e=this.keyW)||void 0===e?void 0:e.isDown)&&(this.body.velocity.y=-s),(null===(n=this.keyA)||void 0===n?void 0:n.isDown)&&(this.body.velocity.x=-s,this.checkFlip(),this.getBody().setOffset(48,15)),(null===(o=this.keyS)||void 0===o?void 0:o.isDown)&&(this.body.velocity.y=s),(null===(i=this.keyD)||void 0===i?void 0:i.isDown)&&(this.body.velocity.x=s,this.checkFlip(),this.getBody().setOffset(15,15)),this.hpValue.setPosition(this.x,this.y-.4*this.height),this.hpValue.setOrigin(.8,.5),this.sendUpdate()},e.prototype.getDamage=function(e){t.prototype.getDamage.call(this,e),this.hpValue.setText(this.hp.toString())},e.prototype.sendUpdate=function(){var t={id:this.playerId,messageType:"status",coordinates:{x:Math.round(this.x),y:Math.round(this.y),directionX:this.scaleX},health:this.hp};JSON.stringify(t)!==JSON.stringify(this.lastSendMessage)&&(this.lastSendMessage=t,window.connection.send(t))},e.prototype.sendAttack=function(){var t={id:this.playerId,messageType:"user_attack"};window.connection.send(t)},e.prototype.initAnimations=function(){this.scene.anims.create({key:"attack",frames:this.scene.anims.generateFrameNames("a-king",{prefix:"attack-",end:2}),frameRate:8})},e}(s.Actor);e.Player=c},729:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Score=e.ScoreOperations=void 0;var s,r=n(67);!function(t){t[t.INCREASE=0]="INCREASE",t[t.DECREASE=1]="DECREASE",t[t.SET_VALUE=2]="SET_VALUE"}(s=e.ScoreOperations||(e.ScoreOperations={}));var a=function(t){function e(e,n,o,i){void 0===i&&(i=0);var s=t.call(this,e,n,o,"Score: ".concat(i))||this;return e.add.existing(s),s.scoreValue=i,s}return i(e,t),e.prototype.changeValue=function(t,e){switch(t){case s.INCREASE:this.scoreValue+=e;break;case s.DECREASE:this.scoreValue-=e;break;case s.SET_VALUE:this.scoreValue=e}this.setText("Score: ".concat(this.scoreValue))},e.prototype.getValue=function(){return this.scoreValue},e}(r.Text);e.Score=a},67:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Text=void 0;var s=function(t){function e(e,n,o,i){var s=t.call(this,e,n,o,i,{fontSize:"calc(100vw / 25)",color:"#fff",stroke:"#000",strokeThickness:4})||this;return s.setOrigin(0,0),e.add.existing(s),s}return i(e,t),e}(n(260).GameObjects.Text);e.Text=s},401:(t,e)=>{Object.defineProperty(e,"__esModule",{value:!0}),e.gameObjectsToObjectPoints=void 0,e.gameObjectsToObjectPoints=function(t){return t.map((function(t){return t}))}},206:(t,e,n)=>{var o=n(260),i=n(580),s=n(800),r=n(595),a=n(321),c={title:"Phaser game",type:Phaser.WEBGL,parent:"game",backgroundColor:"#351f1b",scale:{mode:Phaser.Scale.ScaleModes.NONE,width:window.innerWidth,height:window.innerHeight},physics:{default:"arcade",arcade:{debug:!1}},render:{antialiasGL:!1,pixelArt:!0},callbacks:{postBoot:function(){window.sizeChanged()}},canvasStyle:"display: block; width: 100%; height: 100%;",autoFocus:!0,audio:{disableWebAudio:!1},scene:[i.LoadingScene,s.Level1,r.UIScene]};window.sizeChanged=function(){window.game.isBooted&&setTimeout((function(){window.game.scale.resize(window.innerWidth,window.innerHeight),window.game.canvas.setAttribute("style","display: block; width: ".concat(window.innerWidth,"px; height: ").concat(window.innerHeight,"px;"))}),100)},window.onresize=function(){return window.sizeChanged()},window.connection=new a.Connection,window.game=new o.Game(c)},800:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.Level1=void 0;var s=n(260),r=n(401),a=n(312),c=n(97),l=function(t){function e(){return t.call(this,"level1-scene")||this}return i(e,t),e.prototype.create=function(){this.initMap(),this.player=window.connection.initScene(this),this.physics.add.collider(this.player,this.wallsLayer),this.initEnemies(),this.initCamera()},e.prototype.update=function(){this.player.update()},e.prototype.initMap=function(){this.map=this.make.tilemap({key:"dungeon",tileWidth:16,tileHeight:16}),this.tileset=this.map.addTilesetImage("dungeon","tiles"),this.groundLayer=this.map.createLayer("Ground",this.tileset,0,0),this.wallsLayer=this.map.createLayer("Walls",this.tileset,0,0),this.wallsLayer.setCollisionByProperty({collides:!0}),this.physics.world.setBounds(0,0,this.wallsLayer.width,this.wallsLayer.height)},e.prototype.showDebugWalls=function(){var t=this.add.graphics().setAlpha(.7);this.wallsLayer.renderDebug(t,{tileColor:null,collidingTileColor:new Phaser.Display.Color(243,234,48,255)})},e.prototype.initChests=function(){var t=this,e=(0,r.gameObjectsToObjectPoints)(this.map.filterObjects("Chests",(function(t){return"ChestPoint"===t.name})));this.chests=e.map((function(e){return t.physics.add.sprite(e.x,e.y,"tiles_spr",595).setScale(1.5)})),this.chests.forEach((function(e){t.physics.add.overlap(t.player,e,(function(e,n){t.game.events.emit(a.EVENTS_NAME.chestLoot),n.destroy(),t.cameras.main.flash()}))}))},e.prototype.initCamera=function(){this.cameras.main.setSize(this.game.scale.width,this.game.scale.height),this.cameras.main.startFollow(this.player,!0,.09,.09),this.cameras.main.setZoom(1)},e.prototype.initEnemies=function(){var t=this,e=(0,r.gameObjectsToObjectPoints)(this.map.filterObjects("Enemies",(function(t){return"EnemyPoint"===t.name})));this.enemies=e.map((function(e){return new c.Enemy(t,e.x,e.y,"tiles_spr",t.player,503).setName(e.id.toString()).setScale(1.5)})),this.physics.add.collider(this.enemies,this.wallsLayer),this.physics.add.collider(this.enemies,this.enemies),this.physics.add.collider(this.player,this.enemies,(function(t,e){t.getDamage(1)}))},e}(s.Scene);e.Level1=l},580:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.LoadingScene=void 0;var s=function(t){function e(){return t.call(this,"loading-scene")||this}return i(e,t),e.prototype.preload=function(){window.location.hostname.includes("github.io")?this.load.baseURL="https://maciek1651g.github.io/tanks_js/assets/":this.load.baseURL="./../../assets/",this.load.image("king","sprites/king.png"),this.load.atlas("a-king","spritesheets/a-king.png","spritesheets/a-king_atlas.json"),this.load.image({key:"tiles",url:"tilemaps/tiles/dungeon-16-16.png"}),this.load.tilemapTiledJSON("dungeon","tilemaps/json/dungeon.json"),this.load.spritesheet("tiles_spr","tilemaps/tiles/dungeon-16-16.png",{frameWidth:16,frameHeight:16})},e.prototype.create=function(){this.scene.start("level1-scene"),this.scene.start("ui-scene")},e}(n(260).Scene);e.LoadingScene=s},595:function(t,e,n){var o,i=this&&this.__extends||(o=function(t,e){return o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,e){t.__proto__=e}||function(t,e){for(var n in e)Object.prototype.hasOwnProperty.call(e,n)&&(t[n]=e[n])},o(t,e)},function(t,e){if("function"!=typeof e&&null!==e)throw new TypeError("Class extends value "+String(e)+" is not a constructor or null");function n(){this.constructor=t}o(t,e),t.prototype=null===e?Object.create(e):(n.prototype=e.prototype,new n)});Object.defineProperty(e,"__esModule",{value:!0}),e.UIScene=void 0;var s=n(260),r=n(729),a=n(312),c=function(t){function e(){var e=t.call(this,"ui-scene")||this;return e.chestLootHandler=function(){e.score.changeValue(r.ScoreOperations.INCREASE,10)},e}return i(e,t),e.prototype.create=function(){this.score=new r.Score(this,20,20,0),this.initListeners()},e.prototype.initListeners=function(){this.game.events.on(a.EVENTS_NAME.chestLoot,this.chestLootHandler,this)},e}(s.Scene);e.UIScene=c}},n={};function o(t){var i=n[t];if(void 0!==i)return i.exports;var s=n[t]={exports:{}};return e[t].call(s.exports,s,s.exports,o),s.exports}o.m=e,t=[],o.O=(e,n,i,s)=>{if(!n){var r=1/0;for(h=0;h<t.length;h++){for(var[n,i,s]=t[h],a=!0,c=0;c<n.length;c++)(!1&s||r>=s)&&Object.keys(o.O).every((t=>o.O[t](n[c])))?n.splice(c--,1):(a=!1,s<r&&(r=s));if(a){t.splice(h--,1);var l=i();void 0!==l&&(e=l)}}return e}s=s||0;for(var h=t.length;h>0&&t[h-1][2]>s;h--)t[h]=t[h-1];t[h]=[n,i,s]},o.o=(t,e)=>Object.prototype.hasOwnProperty.call(t,e),(()=>{var t={179:0};o.O.j=e=>0===t[e];var e=(e,n)=>{var i,s,[r,a,c]=n,l=0;if(r.some((e=>0!==t[e]))){for(i in a)o.o(a,i)&&(o.m[i]=a[i]);if(c)var h=c(o)}for(e&&e(n);l<r.length;l++)s=r[l],o.o(t,s)&&t[s]&&t[s][0](),t[s]=0;return o.O(h)},n=self.webpackChunktanks_js=self.webpackChunktanks_js||[];n.forEach(e.bind(null,0)),n.push=e.bind(null,n.push.bind(n))})();var i=o.O(void 0,[426],(()=>o(206)));i=o.O(i)})();