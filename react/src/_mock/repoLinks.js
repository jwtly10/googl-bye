import { faker } from '@faker-js/faker';
import { sample } from 'lodash';

export const generateMockReposLinks = () => {
    const count = faker.number.int({ min: 500, max: 1000 });
    return Array.from({ length: count }, (_, index) => {
        const repoId = index + 1;
        const repoName = faker.helpers
            .uniqueArray(faker.word.words, 2)
            .join('-')
            .toLowerCase()
            .replace(' ', '');
        const repoAuthor = faker.internet.userName();
        return {
            id: faker.string.uuid(),
            avatarUrl: `/assets/images/repos/repo_${repoId}.jpg`,
            name: repoName,
            author: repoAuthor,
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
            parseStatus: sample(['DONE', 'ERROR']),
            apiUrl: `https://api.github.com/repos/${repoAuthor}/${repoName}`,
            ghUrl: `https://github.com/${repoAuthor}/${repoName}`,
            cloneUrl: `https://github.com/${repoAuthor}/${repoName}.git`,
            errorMsg: faker.helpers.maybe(() => faker.lorem.sentence(), { probability: 0.2 }),
            links: generateMockLinks(repoId, repoAuthor, repoName),
        };
    });
};

const generateMockLinks = (repoId, repoAuthor, repoName) => {
    const linkCount = faker.number.int({ min: 0, max: 3 });
    return Array.from({ length: linkCount }, () => {
        const file = faker.system.fileName();
        const lineNumber = faker.number.int({ min: 1, max: 1000 });
        const path = faker.system.filePath();
        return {
            id: faker.string.uuid(),
            repoId: repoId,
            url: `https://goo.gl/${faker.string.alphanumeric(6)}`,
            expandedUrl: faker.internet.url(),
            file: file,
            lineNumber: lineNumber,
            path: path,
            githubUrl: `https://github.com/${repoAuthor}/${repoName}/blob/main/${path}#L${lineNumber}`,
        };
    });
};
