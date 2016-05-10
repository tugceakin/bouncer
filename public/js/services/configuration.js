bouncerApp.factory('configuration', function ($http, $interval) {

   return {
       getBackendServers: function($scope){    

        for(var i=1; i<=$scope.backendServerCount; i++){
            var backendServer = {
              "id" : $scope.configId + "-" + i,
              "host": "backendServer" + i,
              "configId": $scope.configId
            };
            $scope.backendServers.push(backendServer);
            console.log(backendServer);
        }

        return $scope.backendServers;
       }


   }
});