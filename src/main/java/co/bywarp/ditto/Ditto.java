/*
 * Copyright (c) 2019-2020 Warp <legal@warp.pw>
 * All Rights Reserved.
 *
 * This software is proprietary and is designed and intended for internal use only.
 * Unauthorized use, replication, distribution, or modification of this software
 * in any capacity is unlawful and punishable by the full extent of the law.
 */

package co.bywarp.ditto;

import co.bywarp.lightkit.util.AnsiColors;
import co.bywarp.lightkit.util.EnumUtils;
import co.bywarp.lightkit.util.IOUtils;
import co.bywarp.lightkit.util.Pair;
import co.bywarp.lightkit.util.logger.Logger;
import co.bywarp.lightkit.util.timings.Timings;

import org.apache.commons.io.FileUtils;
import org.eclipse.jgit.api.Git;
import org.eclipse.jgit.api.errors.GitAPIException;
import org.eclipse.jgit.lib.Constants;
import org.eclipse.jgit.lib.ObjectId;
import org.eclipse.jgit.transport.UsernamePasswordCredentialsProvider;
import org.json.JSONObject;

import java.io.BufferedWriter;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Calendar;
import java.util.Objects;
import java.util.concurrent.atomic.AtomicReference;
import java.util.regex.Pattern;
import java.util.stream.Stream;

public class Ditto {

    private enum ServerType {
        GAME, LOBBY
    }

    public static void main(String[] args) {
        Logger logger = new Logger("Ditto", Logger.Color.PURPLE);
        logger.raw(AnsiColors.PURPLE + "\n" +
                "\n" +
                "                                           ,--._\n" +
                "                                        _,'     `.\n" +
                "                              ,.-------\"          `.\n" +
                "                             /                 \"    `-.__\n" +
                "                            .         \"        _,        `._\n" +
                "                            |            __..-\"             `.\n" +
                "                            |        ''\"'                     `._\n" +
                "                            |                                    `\"-.\n" +
                "                            '                                        `.\n" +
                "                           .                                          |\n" +
                "                          /                                           |\n" +
                "                       _,'                                           ,'\n" +
                "                     ,\"                                             / \n" +
                "                    .                                              /           ██████╗ ██╗████████╗████████╗ ██████╗  \n" +
                "                    |                                             /            ██╔══██╗██║╚══██╔══╝╚══██╔══╝██╔═══██╗\n" +
                "                    |                                            .             ██║  ██║██║   ██║      ██║   ██║   ██║\n" +
                "                    '                                            |             ██║  ██║██║   ██║      ██║   ██║   ██║\n" +
                "                     `.                                          |             ██████╔╝██║   ██║      ██║   ╚██████╔╝\n" +
                "                       `.                                        |             ╚═════╝ ╚═╝   ╚═╝      ╚═╝    ╚═════╝ \n" +
                "                         `.                                      '\n" +
                "                           .                                      .       " + AnsiColors.RESET + "Initializing " + AnsiColors.PURPLE + "Ditto" + AnsiColors.RESET + " version 0.1/master " + AnsiColors.WHITE + "(snapshot)" + AnsiColors.PURPLE + "\n" +
                "                           |                                       `.         " + AnsiColors.RESET + "Warp (c) 2019-" + Calendar.getInstance().get(Calendar.YEAR) + " ―― internal use only." + AnsiColors.PURPLE + "\n" +
                "                           '                                        |\n" +
                "                         ,'                                         |\n" +
                "                       ,'                                           '\n" +
                "                      /                                _...._      /\n" +
                "                     .                              ,-'      `\"'--'\n" +
                "      ___            |                            ,'\n" +
                "   ,-'   `\"-._     _.'                          ,'\n" +
                "  /           `\"--'             _,....__     _,'\n" +
                " '                            .'        `---'\n" +
                " `                 ____     ,'\n" +
                "  .           _.-'\"    `---'\n" +
                "   `-._    _.\"\n" +
                "       \"\"\"'\n" + AnsiColors.RESET);

        /* Server Properties */
        String serverName = System.getenv("SERVER_NAME");
        ServerType serverType = EnumUtils.matchByName(ServerType.class,
                System.getenv("SERVER_TYPE"));

        /* Git Properties */
        String channel = System.getenv("PLUGIN_CHANNEL");
        String user = System.getenv("GH_USERNAME");
        String token = System.getenv("GH_TOKEN");

        if (channel == null) channel = "production";
        if (serverType == null) {
            logger.severe("Missing server type.");
            System.exit(-1);
            return;
        }

        if (user == null || token == null) {
            logger.severe("Missing login credentials.");
            System.exit(-1);
            return;
        }

        logger.info("Deployment Information");
        logger.info(" - Channel: " + channel);
        logger.info(" - Bean Name: " + serverName);
        logger.info(" - Bean Type: " + serverType.name());

        try {
            Timings timings = new Timings("Ditto", "<init>");
            Git git = Git.cloneRepository()
                    .setURI("https://github.com/bywarp/melon-" + channel)
                    .setCredentialsProvider(
                            new UsernamePasswordCredentialsProvider(user, token)
                    ).call();

            ObjectId id = git
                    .getRepository()
                    .resolve(Constants.HEAD);
            String headSha = id.getName();
            JSONObject versionInfo = new JSONObject()
                    .put("sha", headSha)
                    .put("channel", channel);

            logger.info("[Bundle Info] Release channel: " + channel + ", version: " + headSha);

            File current = new File(Paths.get("").toAbsolutePath().toUri());
            File cloned = new File("melon-" + channel);
            if (!cloned.exists()) {
                logger.severe("Failed to locate cloned directory. (Does it exist?)");
                System.exit(-1);
                return;
            }

            File versionData = new File(current, "version.json");
            BufferedWriter writer = new BufferedWriter(new FileWriter(versionData));
            writer.write(versionInfo.toString(3));
            writer.close();

            logger.info("Upstream VCS data recorded to " + AnsiColors.GREEN + versionData.getAbsolutePath().replaceAll("\\./", "") + AnsiColors.RESET);
            logger.info("Preparing to unpack bundle..");

            File gameParcel = new File(cloned, "game.json");
            File globalParcel = new File(cloned, "parcel.json");
            File clonedPlugins = new File(cloned, "plugins");
            File plugins = new File(current, "plugins");

            if (serverType == ServerType.GAME) {
                FileUtils.moveFileToDirectory(gameParcel, current, false);
            } else {
                FileUtils.forceDelete(gameParcel);
            }

            FileUtils.moveFileToDirectory(globalParcel, current, false);
            FileUtils.copyDirectory(clonedPlugins, plugins);
            logger.info(" - Unpacked preference parcels.");

            File commonLibs = new File(plugins, "common");
            File gameLibs = new File(plugins, "game");
            File lobbyLibs = new File(plugins, "lobby");

            moveAll(commonLibs.listFiles(), plugins, logger, "common libraries");
            FileUtils.deleteDirectory(commonLibs);

            switch (serverType) {
                case GAME:
                    moveAll(gameLibs.listFiles(), plugins, logger, "arcade plugins");
                    break;
                case LOBBY:
                    moveAll(lobbyLibs.listFiles(), plugins, logger, "lobby plugins");
                    break;
            }

            FileUtils.deleteDirectory(gameLibs);
            FileUtils.deleteDirectory(lobbyLibs);
            FileUtils.deleteDirectory(cloned);
            timings.complete("Successfully unpacked bundles in %c%tms%r.");

            /* Configure base server properties */
            rename(new File("server.properties"), logger,
                    Pair.construct(Pattern.compile("server-port=(\\d{5})"), "server-port=" + env("PORT")),
                    Pair.construct(Pattern.compile("server-name=SERVER"), "server-name=" + env("SERVER_NAME")),
                    Pair.construct(Pattern.compile("online-mode=false"), "online-mode=" + env("ONLINE_MODE")));

            /* Configure bungeecord mode */
            rename(new File("spigot.yml"), logger,
                    Pair.construct(Pattern.compile("bungeecord: (true|false)"), "bungeecord: " + env("BUNGEE")));

            /* Configure global preference parcel */
            rename(globalParcel, logger,
                    Pair.construct(Pattern.compile("DEFAULT_SERVER_NAME"), env("SERVER_NAME")),
                    Pair.construct(Pattern.compile("DEFAULT_SERVER_GROUP"), env("SERVER_GROUP")),
                    Pair.construct(Pattern.compile("DEFAULT_SERVER_TYPE"), env("SERVER_TYPE")),
                    Pair.construct(Pattern.compile("DEFAULT_SERVER_REGION"), env("SERVER_REGION")),
                    Pair.construct(Pattern.compile("DEFAULT_JOIN_QUIT"), env("JOIN_QUIT")),
                    Pair.construct(Pattern.compile("DEFAULT_INCOGNITO"), env("INCOGNITO")),
                    Pair.construct(Pattern.compile("DEFAULT_DATABASE_NAME"), env("DATABASE_NAME")),
                    Pair.construct(Pattern.compile("DEFAULT_DATABASE_USERNAME"), env("DATABASE_USERNAME")),
                    Pair.construct(Pattern.compile("DEFAULT_DATABASE_PASSWORD"), env("DATABASE_PASSWORD")),
                    Pair.construct(Pattern.compile("DEFAULT_DATABASE_HOST"), env("DATABASE_HOST")),
                    Pair.construct(Pattern.compile("DEFAULT_GATEWAY_HOST"), env("STARGATE_HOST")),
                    Pair.construct(Pattern.compile("DEFAULT_GATEWAY_ID"), env("STARGATE_CLIENT_ID")),
                    Pair.construct(Pattern.compile("DEFAULT_GATEWAY_SECRET"), env("STARGATE_CLIENT_SECRET")),
                    Pair.construct(Pattern.compile("DEFAULT_REDIS_HOST"), env("REDIS_HOST")),
                    Pair.construct(Pattern.compile("DEFAULT_REDIS_PORT"), env("REDIS_PORT")),
                    Pair.construct(Pattern.compile("DEFAULT_REDIS_AUTH"), env("REDIS_AUTH")),
                    Pair.construct(Pattern.compile("DEFAULT_REDIS_PASSWORD"), env("REDIS_PASSWORD")),
                    Pair.construct(Pattern.compile("DEFAULT_REDIS_SPEAKER"), env("REDIS_SPEAKER")));

            /* If applicable, configure game preference parcel */
            if (gameParcel.exists()) {
                rename(gameParcel, logger,
                        Pair.construct(Pattern.compile("DEFAULT_GAME_NAME"), env("GAME")),
                        Pair.construct(Pattern.compile("DEFAULT_GAME_ROTATE"), env("GAME_ROTATE")),
                        Pair.construct(Pattern.compile("DEFAULT_MAX_PLAYERS"), env("MAX_PLAYERS")),
                        Pair.construct(Pattern.compile("DEFAULT_AWARD_STATS"), env("AWARD_STATS")),
                        Pair.construct(Pattern.compile("DEFAULT_HOST_UUID"), env("GAME_SERVER_HOST")),
                        Pair.construct(Pattern.compile("DEFAULT_IS_HOSTED"), env("GAME_SERVER_IS_HOSTED")),
                        Pair.construct(Pattern.compile("DEFAULT_GAME_SERVER_TYPE"), env("GAME_SERVER_TYPE")),
                        Pair.construct(Pattern.compile("DEFAULT_GAME_SERVER_PRIORITY"), env("GAME_SERVER_PRIORITY")),
                        Pair.construct(Pattern.compile("DEFAULT_IS_EVENT"), env("EVENT")));
            }

            logger.info("Every day is a good day when you're playing on " + AnsiColors.GREEN + "Melon Games" + AnsiColors.RESET + "!");
        } catch (GitAPIException e) {
            logger.except(e, "Exception cloning " + channel + " bundle");
            System.exit(-1);
        } catch (IOException e) {
            logger.except(e, "Exception unpacking " + channel + " bundle");
            System.exit(-1);
        }
    }

    private static void moveAll(File[] files, File destination, Logger logger, String what) {
        Arrays.stream(Objects.requireNonNull(files)).forEach(file -> {
            try {
                FileUtils.moveFileToDirectory(file, destination, false);
            } catch (IOException e) {
                logger.except(e, "Failed to move binaries");
            }
        });

        logger.info(" - Unpacked " + what + ".");
    }

    @SafeVarargs
    private static void rename(File file, Logger logger, Pair<Pattern, String>... entries) {
        try {
            Timings timings = new Timings("Ditto", "<init>");
            String rawContents = IOUtils.toString(file);
            if (rawContents.isEmpty()) {
                throw new IOException("Empty or malformed data");
            }

            AtomicReference<String> contents = new AtomicReference<>(rawContents);
            Stream.of(entries).forEach(ent -> {
                Pattern regex = ent.getK();
                String replacer = ent.getV();

                contents.getAndUpdate(s -> s.replaceAll(regex.pattern(), replacer));
            });

            BufferedWriter writer = new BufferedWriter(new FileWriter(file));
            writer.write(contents.get());
            writer.close();

            timings.complete("Configured " + file.getName() + " in %c%tms%r.");
        } catch (IOException e) {
            logger.except(e, "Fatal exception configuring " + file.getName());
            System.exit(-1);
        }
    }

    private static String env(String key) {
        return System.getenv(key);
    }

}
