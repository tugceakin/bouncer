bouncerApp.controller('ConfigurationController', function($scope, $http, $window, benchmarking, configuration) {

    $scope.backendServerCount = 1;

    $scope.addConfiguration = function(config){
        var backendServerCount = $window.document.getElementsByClassName("backendServerInput").length;
        configuration.addConfiguration($scope, config, backendServerCount);
    };

    $scope.updateConfiguration = function(){

    };

    $scope.removeConfiguration = function(){

    };

});