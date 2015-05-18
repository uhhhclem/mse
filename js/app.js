var mseApp = angular.module('mseApp', ['ngMaterial']);

mseApp.controller('mainCtrl', function($scope, $http){
    $http.get('/api/newGame').success(function(d){
        $scope.board = d;
    });
})