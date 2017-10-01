# Troubleshooting

## Failing to run the container, permission issue

This may be due to your Docker installation that is not authorized
to be used by a non-root user.

You need to add your user to the `docker` group (it should have been
created during the installation).

This this [Stackoverflow question](https://stackoverflow.com/questions/35849533/running-docker-as-non-root-user)
for more information.

