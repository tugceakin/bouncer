bouncerApp.controller('ConfigurationController', function($scope, $http, $window, benchmarking, configuration) {

    $scope.backendServerCount = 1;
    configuration.getAllConfigs($scope);


    $scope.addConfiguration = function(config){
        var backendServerCount = $window.document.getElementsByClassName("backendServerInput").length;
        configuration.addConfiguration($scope, config, backendServerCount);
        configuration.getAllConfigs($scope);
    };

    $scope.updateConfiguration = function(config){
    };

    $scope.removeConfiguration = function(config){
        configuration.removeConfiguration($scope, config);
    };

});