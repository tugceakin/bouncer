bouncerApp.factory('configuration', function ($http, $interval) {


  var getBackendServers = function($scope,config, backendServerCount){
      var backendServers = [];

    for(var i=1; i<=backendServerCount; i++){
          var serverHostName = "backendServer" + i;
          var backendServer = {
            //"id" : config.configId + "" + i,
            "host": config[serverHostName]
            //"configId": config.configId
          };
          backendServers.push(backendServer);
          console.log(backendServer);
      }
      console.log(backendServers.length);
      return backendServers;
  };

   return {

       addConfiguration: function($scope, config, backendServerCount){

        var url = "http://localhost:8080/addConfiguration";
        var config = {
            //"id": config.configId,
            "host": config.hostName,
            "path": config.path,
            "targetPath": config.targetPath,
            "reqPerSecond": config.reqPerSecond,
            "concurrency": config.concurrency,
            "backendServers": getBackendServers($scope, config, backendServerCount)
        };

        console.log(config);

        $http({
            method: 'POST',
            url: url,
            headers: {'Content-Type': 'application/json'}, 
            data: JSON.stringify(config)
        }).success(function (data) {
            console.log(JSON.stringify(data));
        });
       }
   }
});