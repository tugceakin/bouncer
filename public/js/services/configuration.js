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
      return backendServers;
  };

   return {

       addConfiguration: function($scope, config, backendServerCount){
        var url = "http://localhost:8080/addConfiguration";
        if(config.targetPath == undefined) config.targetPath = ""
        if(config.path == undefined) config.path = ""
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
       },

       removeConfiguration: function($scope, config){
          var url = "http://localhost:8080/removeConfiguration";
          $http({
              method: 'POST',
              url: url,
              headers: {'Content-Type': 'application/json'}, 
              data: JSON.stringify(config)
          }).success(function (data) {
              console.log(JSON.stringify(data));
          });
       },

       getAllConfigs: function($scope){
        var url = "http://localhost:8080/getAllConfigs";
        $http({
            method: 'POST',
            url: url,
            headers: {'Content-Type': 'application/json'}
        }).success(function (data) {
            console.log(JSON.stringify(data));
            $scope.configs = data;
        });
       }
   }
});