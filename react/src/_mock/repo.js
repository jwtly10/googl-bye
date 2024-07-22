import { faker } from '@faker-js/faker';
import { sample } from 'lodash';

export const generateMockRepos = () => {
    // Generate a random number between 500 and 1000
    const count = faker.number.int({ min: 500, max: 1000 });

    return Array.from({ length: count }, (_, index) => ({
        id: faker.string.uuid(),
        avatarUrl: `/assets/images/repos/repo_${index + 1}.jpg`,
        name: faker.helpers.uniqueArray(faker.word.words, 2).join('-').toLowerCase().replace(' ', ''),
        author: faker.internet.userName(),
        isVerified: faker.datatype.boolean(),
        language: sample([
            'JavaScript',
            'Python',
            'Java',
            'C++',
            'Ruby',
            'Go',
            'TypeScript',
            'PHP',
            'C#',
            'Swift',
        ]),
        stars: faker.number.int({ min: 0, max: 10000 }),
        forks: faker.number.int({ min: 0, max: 5000 }),
        lastCommit: faker.date.recent({ days: 30 }),
        size: faker.number.int({ min: 100, max: 1000000 }), // size in KB
        parseStatus: sample(['PENDING', 'PROCESSING', 'DONE', 'ERROR']),
        apiUrl: `https://api.github.com/repos/${faker.internet.userName()}/${faker.helpers
            .uniqueArray(faker.word.words, 2)
            .join('-')
            .toLowerCase()}`,
        ghUrl: `https://github.com/${faker.internet.userName()}/${faker.helpers
            .uniqueArray(faker.word.words, 2)
            .join('-')
            .toLowerCase()}`,
        cloneUrl: `https://github.com/${faker.internet.userName()}/${faker.helpers
            .uniqueArray(faker.word.words, 2)
            .join('-')
            .toLowerCase()}.git`,
        errorMsg: faker.helpers.maybe(() => faker.lorem.sentence(), { probability: 0.2 }),
    }));
};
